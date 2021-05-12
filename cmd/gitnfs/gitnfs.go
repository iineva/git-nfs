package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/iineva/git-nfs/pkg/gitnfs"
	"github.com/iineva/git-nfs/pkg/logger"
	"github.com/iineva/git-nfs/pkg/signal"
)

// cli args
type Args struct {
	Help  bool
	Debug bool

	SyncInterval time.Duration

	// nfs
	Addr     string
	Readonly bool

	// git
	GitURL           string
	GitCommitName    string
	GitCommitEmail   string
	GitReferenceName string
}

func parseArgs(args *Args) error {

	flag.Usage = usage

	flag.BoolVar(&args.Help, "h", false, "this help")
	flag.BoolVar(&args.Debug, "d", false, "enable debug logs")

	flag.StringVar(&args.Addr, "a", ":0", "nfs listen addr")
	flag.BoolVar(&args.Readonly, "o", false, "make nfs server readonly")

	flag.StringVar(&args.GitCommitName, "m", "gitnfs", "git commit name")
	flag.StringVar(&args.GitCommitEmail, "e", "gitnfs@example.com", "git commit email")
	flag.StringVar(&args.GitReferenceName, "r", "refs/heads/main", "git reference name")

	flag.DurationVar(&args.SyncInterval, "s", time.Second*5, "interval when sync nfs files to git repo")

	flag.Parse()
	if len(flag.Args()) != 1 {
		return errors.New("you need help")
	}
	args.GitURL = flag.Args()[0]

	return nil
}

func main() {

	args := new(Args)
	if err := parseArgs(args); err != nil || args.Help {
		usage()
		os.Exit(0)
	}

	if args.GitURL == "" {
		flag.Usage()
		os.Exit(0)
	}

	// TODO:
	if args.Readonly {
		panic(errors.New("readonly mode is not yet implement"))
	}

	// save memory
	debug.SetGCPercent(1)

	// setup log
	logger.RedirectStd()
	logger.Debug(args.Debug)
	log := logger.New("main")
	log.Debugf("args: %+v", args)

	log.Infof("server starting addr %v", args.Addr)
	gn := gitnfs.New(gitnfs.Config{
		Addr:             args.Addr,
		GitURL:           args.GitURL,
		GitCommitName:    args.GitCommitName,
		GitCommitEmail:   args.GitCommitEmail,
		GitReferenceName: args.GitReferenceName,
		SyncInterval:     args.SyncInterval,
	})
	signal.AddTermCallback(func(s os.Signal, done func()) {
		gn.Close()
		done()
	})
	err := gn.Serve()
	if err != nil {
		log.Error(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `gitnfs version: 0.1.0
Usage: gitnfs [options] https://github.com/iineva/gitnfs

Options:
`)
	flag.PrintDefaults()
}
