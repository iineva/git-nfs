package git

import (
	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (g *Git) Clone() error {
	// clone
	r, err := go_git.Clone(memory.NewStorage(), g.config.FS, &go_git.CloneOptions{
		URL:           g.config.URL,
		Depth:         g.depth(),
		ReferenceName: g.referenceName(),
		Auth:          g.config.Auth,
		// Progress:      os.Stdout,
	})
	if err != nil {
		return err
	}
	g.repo = r
	return nil
}
