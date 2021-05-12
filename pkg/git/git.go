package git

import (
	"github.com/go-git/go-billy/v5"
	go_git "github.com/go-git/go-git/v5"
)

type Config struct {
	URL           string           // required
	ReferenceName string           // option, default: master
	RemoteName    string           // option, default: origin
	Depth         int              // option, default: 1
	CommitName    string           // required
	CommitEmail   string           // required
	FS            billy.Filesystem // required
}

type Git struct {
	config Config
	repo   *go_git.Repository
}

func New(conf Config) *Git {
	g := &Git{config: conf}
	return g
}

func (g *Git) Close() error {
	g.repo = nil
	return nil
}
