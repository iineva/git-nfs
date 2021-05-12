package nfs

import (
	"net"

	"github.com/willscott/go-nfs"
	"github.com/willscott/go-nfs/filesystem"
	"github.com/willscott/go-nfs/helpers"
)

type NFS struct {
	listener net.Listener
	fs       filesystem.FS
}

func New(listener net.Listener, fs filesystem.FS) *NFS {
	return &NFS{listener: listener, fs: fs}
}

func (n *NFS) Serve() error {
	handler := helpers.NewNullAuthHandler(n.fs)
	cacheHelper := helpers.NewCachingHandler(handler, 1024)
	return nfs.Serve(n.listener, cacheHelper)
}

func (n *NFS) Close() error {
	return n.listener.Close()
}
