package gitnfs

import "testing"

func Test_gitnfs(t *testing.T) {
	conf := Config{
		Addr:           ":5566",
		GitURL:         "/Users/steven/test/git-nfs-test",
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
	err := gn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
