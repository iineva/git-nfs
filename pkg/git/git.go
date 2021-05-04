package git

import (
	"os"

	go_git "github.com/go-git/go-git/v5"
)

func Clone() {
	_, err := go_git.PlainClone("/tmp/foo", false, &go_git.CloneOptions{
			URL:      "https://github.com/iineva/templates",
			Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}


}