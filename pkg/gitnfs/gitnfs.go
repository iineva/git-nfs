package gitnfs

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
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
	GitURL           string
	GitCommitName    string
	GitCommitEmail   string
	GitReferenceName string
	GitAuth          transport.AuthMethod
	SyncInterval     time.Duration
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

	// init git
	gn.gitFS = memfs.New()
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
	gn.nfsFS = basefs.NewMemMapFS()
	gn.nfs = nfs.New(nfsListener, gn.nfsFS)

	// sync file to nfs
	if err := syncfile.Billy2Afero(gn.gitFS, gn.nfsFS); err != nil {
		gn.logger.Error("sync git to nfs error:", err)
		return err
	}

	go gn.syncLoop()

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
		if err := syncfile.Afero2Billy(gn.nfsFS, gn.gitFS); err != nil {
			gn.logger.Error("sync error: ", err)
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
