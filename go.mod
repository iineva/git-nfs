module github.com/iineva/git-nfs

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/go-git/go-billy/v5 v5.1.0
	github.com/go-git/go-git/v5 v5.3.0
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/spf13/afero v1.6.0
	github.com/uber/jaeger-client-go v2.28.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/willscott/go-nfs v0.0.0-20210422193315-8a05dee1ebbe
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
)

replace (
	// afero not support get uid and gid for now, use this fork to support
	github.com/spf13/afero => github.com/iineva/afero v1.6.1-0.20210510115905-57c673cfea7b
	github.com/willscott/go-nfs => github.com/iineva/go-nfs v0.0.0-20210512034119-3d40f31ee9e6
)
