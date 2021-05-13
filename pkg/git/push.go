package git

import (
	go_git "github.com/go-git/go-git/v5"
)

func (g *Git) Push() error {
	return g.repo.Push(&go_git.PushOptions{
		RemoteName: g.remoteName(),
		Auth:       g.config.Auth,
		// Progress: os.Stdout,
	})
}
