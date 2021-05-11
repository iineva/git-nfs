package main

import (
	"log"
	"os"
	"runtime/debug"

	git_nfs "github.com/iineva/git-nfs/pkg/git"
)

// hook os.Stdout output from go-nfs module
func init() {
	f, err := os.OpenFile("/tmp/go-nfs.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	os.Stdout = f
	os.Stderr = f
	log.SetOutput(f)
}

func main() {
	debug.SetGCPercent(1)
	git_nfs.Start()
}
