package gitnfs

import (
	"errors"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/afero"
	"github.com/willscott/go-nfs/filesystem"
	"github.com/willscott/go-nfs/filesystem/basefs"
	"go.uber.org/zap"

	"github.com/iineva/git-nfs/pkg/git"
	"github.com/iineva/git-nfs/pkg/logger"
	"github.com/iineva/git-nfs/pkg/nfs"
	"github.com/iineva/git-nfs/pkg/syncfile"
)

type Config struct {
	Addr             string
	GitURL           *url.URL
	GitCommitName    string
	GitCommitEmail   string
	GitReferenceName string
	GitAuth          transport.AuthMethod
	SyncInterval     time.Duration
	CacheDir         string
	Readonly         bool
}

type GitNFS struct {
	config Config
	logger *zap.SugaredLogger

	git    *git.Git
	nfs    *nfs.NFS
	gitFS  billy.Filesystem
	nfsFS  filesystem.FS
	runing int32
}

func New(conf Config) *GitNFS {
	return &GitNFS{
		config: conf,
		logger: logger.New("gitnfs"),
	}
}

func (gn *GitNFS) Serve() error {

	if gn.hasCacheDir() {
		if f, err := os.Stat(gn.config.CacheDir); os.IsNotExist(err) {
			if err := os.MkdirAll(gn.config.CacheDir, 0755); err != nil {
				return err
			}
		} else {
			if !f.IsDir() {
				return errors.New("cache path is not a dir")
			}
		}
	}

	// init git
	var gitFS billy.Filesystem = nil
	if gn.hasCacheDir() {
		gitFS = osfs.New(gn.config.CacheDir)
	} else {
		gitFS = memfs.New()
	}
	gn.gitFS = gitFS
	gn.git = git.New(git.Config{
		FS:            gn.gitFS,
		URL:           gn.config.GitURL,
		Auth:          gn.config.GitAuth,
		CommitName:    gn.config.GitCommitName,
		CommitEmail:   gn.config.GitCommitEmail,
		ReferenceName: gn.config.GitReferenceName,
	})
	err := gn.git.Clone()
	if err != nil {
		switch err {
		case transport.ErrEmptyRemoteRepository:
			gn.logger.Debug("clone: ", err)
		case transport.ErrAuthenticationRequired, transport.ErrAuthorizationFailed, transport.ErrInvalidAuthMethod:
			gn.logger.Error("clone: ", err)
			return err
		default:
			gn.logger.Error("clone: ", err)
			return err
		}
	}

	// init nfs
	nfsListener, err := net.Listen("tcp", gn.config.Addr)
	gn.logger.Info("listening port ", nfsListener.Addr().(*net.TCPAddr).Port)
	if err != nil {
		gn.logger.Errorf("listen addr: %v %v", gn.config.Addr, err)
		return err
	}
	var nfsFS filesystem.FS = nil
	if gn.hasCacheDir() {
		nfsFS = basefs.NewOsFS(filepath.Join(gn.config.CacheDir, git.CheckoutDir))
	} else {
		nfsFS = basefs.NewMemMapFS()
	}
	if gn.config.Readonly {
		if fs, ok := nfsFS.(basefs.BaseFS); ok {
			fs.SetSource(afero.NewReadOnlyFs(fs.GetSource()))
		}
	}
	gn.nfsFS = nfsFS
	gn.nfs = nfs.New(nfsListener, gn.nfsFS)

	// in memory mode, sync file to nfs
	// TODO: avoid to copy
	if !gn.hasCacheDir() {
		if err := syncfile.Billy2Afero(gn.gitFS, gn.nfsFS); err != nil {
			gn.logger.Error("sync git to nfs error:", err)
			return err
		}
	}

	if !gn.config.Readonly {
		go gn.syncLoop()
	}

	return gn.nfs.Serve()
}

// sync file --> pull --> push
func (gn *GitNFS) syncLoop() {
	atomic.AddInt32(&gn.runing, 1)
	duration := gn.config.SyncInterval
	if duration <= 0 {
		duration = time.Second * 5
	}
	for gn.runing > 0 {
		time.Sleep(duration)
		if gn.runing <= 0 {
			break
		}
		// in memory mode, sync file
		// TODO: avoid to copy
		if !gn.hasCacheDir() {
			if err := syncfile.Afero2Billy(gn.nfsFS, gn.gitFS); err != nil {
				gn.logger.Error("sync error: ", err)
			}
		}
		if err := gn.git.Commit(); err != nil {
			gn.logger.Error("commit: ", err)
		}
		if err := gn.git.Pull(); err != nil {
			if err == go_git.NoErrAlreadyUpToDate {
				gn.logger.Debug("pull: ", err)
			} else {
				gn.logger.Error("pull: ", err)
			}
		}
		if err := gn.git.Push(); err != nil {
			if err == go_git.NoErrAlreadyUpToDate {
				gn.logger.Debug("push: ", err)
			} else {
				gn.logger.Error("push: ", err)
			}
		}
	}
}

func (gn *GitNFS) Close() error {
	atomic.AddInt32(&gn.runing, 0)
	if err := gn.nfs.Close(); err != nil {
		return err
	}
	if err := gn.git.Close(); err != nil {
		return err
	}
	return nil
}

func (gn *GitNFS) hasCacheDir() bool {
	return gn.config.CacheDir != ""
}
