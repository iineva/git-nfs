package git

import "github.com/go-git/go-git/v5/plumbing"

func (g *Git) remoteName() string {
	remoteName := g.config.RemoteName
	if remoteName == "" {
		remoteName = "origin"
	}
	return remoteName
}

func (g *Git) referenceName() plumbing.ReferenceName {
	referenceName := g.config.ReferenceName
	if referenceName == "" {
		referenceName = "main"
	}
	return plumbing.ReferenceName(referenceName)
}

func (g *Git) depth() int {
	depth := g.config.Depth
	if depth == 0 {
		depth = 1
	}
	return depth
}
