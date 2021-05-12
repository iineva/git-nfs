package git

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/iineva/git-nfs/pkg/afero2billy"
	nfs "github.com/willscott/go-nfs"
	"github.com/willscott/go-nfs/filesystem/basefs"
	nfshelper "github.com/willscott/go-nfs/helpers"

	app_logger "github.com/iineva/git-nfs/pkg/logger"
)

var logger *zap.SugaredLogger

func init() {
	logger = app_logger.New("git-nfs")
	// TODO: sync when exit
	// defer logger.Sync() // flushes buffer, if any
}

func PushGit(r *go_git.Repository) error {
	// return nil
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// check file changes
	s, err := w.Status()
	if err != nil {
		return err
	}
	if s.IsClean() {
		logger.Debug("======= all files clean =======")
		return nil
	}

	// add file changes
	if err := w.AddWithOptions(&go_git.AddOptions{
		All: true,
	}); err != nil {
		return err
	}

	// commit
	commit, err := w.Commit(fmt.Sprintf("Update at %v", time.Now()), &go_git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "Steven Jobs",
			Email: "s@ineva.cn",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	if _, err := r.CommitObject(commit); err != nil {
		return err
	}

	if err := r.Push(&go_git.PushOptions{
		// RemoteName: "main",
		// Progress: os.Stdout,
	}); err != nil {
		return err
	}
	logger.Info("======= push done =======")

	return nil
}

func Start() {

	tracerCloser := StarTracer()
	// TODO: close when exit
	defer tracerCloser.Close()

	// init git
	memFs := memfs.New()
	cloneSpan, _ := opentracing.StartSpanFromContext(context.Background(), "git-clone")
	r, err := go_git.Clone(memory.NewStorage(), memFs, &go_git.CloneOptions{
		URL: "/Users/steven/test/git-nfs-test",
		// Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		Depth:         1, // if memery mode
	})
	cloneSpan.Finish()
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		logger.Errorf("clone error: %v", err.Error())
		panic(err)
	}

	// start nfs
	listener, err := net.Listen("tcp", ":5566")
	if err != nil {
		logger.Errorf("Failed to listen: %v\n", err)
		return
	}
	logger.Infof("Server running at %s", listener.Addr())
	memMapFS := basefs.NewMemMapFS()
	handler := nfshelper.NewNullAuthHandler(memMapFS)
	cacheHelper := nfshelper.NewCachingHandler(handler, 1024)

	// init sync file
	syncBillySpan, _ := opentracing.StartSpanFromContext(context.Background(), "sync-billy-to-afero")
	if err := afero2billy.SyncBilly2Afero(memFs, memMapFS); err != nil {
		panic(err)
	}
	syncBillySpan.Finish()

	// sync file and push
	const SYNC_INTERVAL = time.Second * 5
	go (func() {
		for {
			time.Sleep(SYNC_INTERVAL)

			syncToGit, ctx := opentracing.StartSpanFromContext(context.Background(), "sync-to-git")

			syncAferoToBillySpan, _ := opentracing.StartSpanFromContext(ctx, "sync-afero-to-billy")
			if err := afero2billy.SyncAfero2Billy(memMapFS, memFs); err != nil {
				logger.Errorf("sync error: %s", err.Error())
			}
			syncAferoToBillySpan.Finish()
			pushSpan, _ := opentracing.StartSpanFromContext(ctx, "git-push")
			if err := PushGit(r); err != nil {
				logger.Errorf("push error: %s", err.Error())
			}
			pushSpan.Finish()

			syncToGit.Finish()
		}
	})()

	// starting
	fmt.Printf("%v", nfs.Serve(listener, cacheHelper))
}
