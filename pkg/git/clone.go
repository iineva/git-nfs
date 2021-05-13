package git

import (
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

const (
	DirPerm     = 0755
	StorerDir   = "storer"
	CheckoutDir = "checkout"
)

func (g *Git) Clone() error {
	// init dirs
	if err := g.config.FS.MkdirAll(StorerDir, DirPerm); err != nil {
		return err
	}
	if err := g.config.FS.MkdirAll(CheckoutDir, DirPerm); err != nil {
		return err
	}

	// clone
	storerFS := osfs.New(filepath.Join(g.config.FS.Root(), StorerDir))
	checkoutFS := osfs.New(filepath.Join(g.config.FS.Root(), CheckoutDir))
	r, err := go_git.Clone(
		filesystem.NewStorage(storerFS, cache.NewObjectLRUDefault()),
		checkoutFS,
		&go_git.CloneOptions{
			URL: g.config.URL.String(),
			// TODO: fix this
			// Depth:         g.depth(),
			ReferenceName: g.referenceName(),
			SingleBranch:  true,
			Auth:          g.config.Auth,
			// Progress:      os.Stdout,
		},
	)
	if err != nil {
		return err
	}
	g.repo = r
	return nil
}
