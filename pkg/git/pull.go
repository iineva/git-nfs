package git

import (
	go_git "github.com/go-git/go-git/v5"
)

func (g *Git) Pull() error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	return w.Pull(&go_git.PullOptions{
		Depth:         g.depth(),
		RemoteName:    g.remoteName(),
		SingleBranch:  true,
		ReferenceName: g.referenceName(),
		Auth:          g.config.Auth,
	})
}
