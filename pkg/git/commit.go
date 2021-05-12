package git

import (
	"fmt"
	"time"

	go_git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (g *Git) Commit() error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	// check file changes
	s, err := w.Status()
	if err != nil {
		return err
	}
	if s.IsClean() {
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
			Name:  g.config.CommitName,
			Email: g.config.CommitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	if _, err := g.repo.CommitObject(commit); err != nil {
		return err
	}

	return nil
}
