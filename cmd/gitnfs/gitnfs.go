package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/iineva/git-nfs/pkg/git"
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
	Addr           string
	Readonly       bool
	CacheDir       string
	CacheDirDelete bool

	// git
	GitURL           string
	GitCommitName    string
	GitCommitEmail   string
	GitReferenceName string

	// Auth
	GitUsername       string
	GitPassword       string
	GitSSHKey         string
	GitSSHKeyFile     string
	GitSSHKeyPassword string
}

func parseArgs(args *Args) error {

	flag.Usage = usage

	flag.BoolVar(&args.Help, "h", false, "this help")
	flag.BoolVar(&args.Debug, "d", false, "enable debug logs")
	flag.StringVar(&args.CacheDir, "c", "", "file cache dir, if it's empty will storage files in memery")
	flag.BoolVar(&args.CacheDirDelete, "z", false, "delete file cache dir when exit")

	flag.StringVar(&args.Addr, "a", ":0", "nfs listen addr")
	flag.BoolVar(&args.Readonly, "o", false, "make nfs server readonly")

	flag.StringVar(&args.GitCommitName, "m", "gitnfs", "git commit name")
	flag.StringVar(&args.GitCommitEmail, "e", "gitnfs@example.com", "git commit email")
	flag.StringVar(&args.GitReferenceName, "r", "refs/heads/main", "git reference name")
	flag.StringVar(&args.GitUsername, "u", "", "basic auth user name")
	flag.StringVar(&args.GitPassword, "p", "", "basic auth password or GitHub personal access token")
	flag.StringVar(&args.GitSSHKey, "k", "", "private key string")
	flag.StringVar(&args.GitSSHKeyFile, "f", "", "private key file")
	flag.StringVar(&args.GitSSHKeyPassword, "K", "", "private key password")

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

	// save memory
	if args.CacheDir == "" {
		debug.SetGCPercent(1)
	}

	// setup log
	logger.RedirectStd("/dev/null")
	logger.Debug(args.Debug)
	log := logger.New("main")
	log.Debugf("args: %+v", args)

	var auth transport.AuthMethod = nil
	uri, err := git.Parse(args.GitURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	if args.GitSSHKey != "" {
		publicKeys, err := ssh.NewPublicKeys(uri.User.Username(), []byte(args.GitSSHKey), args.GitSSHKeyPassword)
		if err != nil {
			log.Fatal("ssh private key file ", err)
			os.Exit(-1)
		}
		log.Debug(publicKeys)
		auth = publicKeys
	} else if args.GitSSHKeyFile != "" {
		publicKeys, err := ssh.NewPublicKeysFromFile(uri.User.Username(), args.GitSSHKeyFile, args.GitSSHKeyPassword)
		if err != nil {
			log.Fatal("ssh private key file ", err)
			os.Exit(-1)
		}
		auth = publicKeys
	} else if args.GitUsername != "" {
		auth = &http.BasicAuth{
			Username: args.GitUsername,
			Password: args.GitPassword,
		}
	}

	log.Infof("server starting addr %v", args.Addr)
	gn := gitnfs.New(gitnfs.Config{
		Addr:             args.Addr,
		GitURL:           uri,
		GitAuth:          auth,
		GitCommitName:    args.GitCommitName,
		GitCommitEmail:   args.GitCommitEmail,
		GitReferenceName: args.GitReferenceName,
		SyncInterval:     args.SyncInterval,
		CacheDir:         args.CacheDir,
		Readonly:         args.Readonly,
	})
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	signal.AddTermCallback(func(s os.Signal, done func()) {
		gn.Close()
		if args.CacheDir != "" && args.CacheDirDelete {
			if err := os.RemoveAll(args.CacheDir); err != nil {
				log.Error(err)
			} else {
				log.Info("cache dir is deleted")
			}
		}
		waitGroup.Done()
		done()
	})
	err = gn.Serve()
	if err != nil {
		log.Error(err)
	}
	// wait until cleanup
	waitGroup.Wait()
}

func usage() {
	fmt.Fprintf(os.Stderr, `gitnfs version: 0.1.0
Usage: gitnfs [options] <YOUR_GIT_REPO_URL>

Options:
`)
	flag.PrintDefaults()
}
