package gitnfs

import (
	"testing"

	"github.com/iineva/git-nfs/pkg/git"
)

func Test_gitnfs(t *testing.T) {
	uri, err := git.Parse("/Users/steven/test/git-nfs-test")
	if err != nil {
		t.Fatal(err)
	}
	conf := Config{
		Addr:           ":5566",
		GitURL:         uri,
		GitCommitName:  "Steven",
		GitCommitEmail: "s@ineva.cn",
	}
	gn := New(conf)
	go func() {
		err := gn.Serve()
		if err != nil {
			t.Fatal(err)
		}
	}()
	if err := gn.Close(); err != nil {
		t.Fatal(err)
	}
}
