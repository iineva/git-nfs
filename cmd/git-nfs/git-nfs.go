package main

import (
	"os"
	"runtime/debug"

	git_nfs "github.com/iineva/git-nfs/pkg/git"
	"github.com/iineva/git-nfs/pkg/logger"
	"github.com/iineva/git-nfs/pkg/signal"
)

func main() {

	logger.RedirectStd()

	log := logger.New("main")
	log.Info("server start!")
	signal.AddTermCallback(func(s os.Signal, done func()) {
		// TODO: close
		done()
	})
	log.Info("server start done!")

	debug.SetGCPercent(1)
	git_nfs.Start()
}
