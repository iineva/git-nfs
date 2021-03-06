module github.com/iineva/git-nfs

go 1.16

require (
	github.com/go-git/go-billy/v5 v5.3.1
	// github.com/go-git/go-git/v5 v5.3.0
	github.com/go-git/go-git/v5 v5.3.1-0.20210512204207-720c192831a8
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/spf13/afero v1.6.0
	github.com/whilp/git-urls v1.0.0
	github.com/willscott/go-nfs v0.0.0-20210422193315-8a05dee1ebbe
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
)

replace (
	// afero not support get uid and gid for now, use this fork to support
	github.com/spf13/afero => github.com/iineva/afero v1.6.1-0.20210510115905-57c673cfea7b
	github.com/willscott/go-nfs => github.com/iineva/go-nfs v0.0.0-20210513170604-595d4471e803
)
