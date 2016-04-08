package fs;

import (
//	"github.com/hanwen/go-fuse/fuse"
//	"github.com/hanwen/go-fuse/fuse/pathfs"
//	"github.com/hanwen/go-fuse/fuse/nodefs"
);


// path rewriting FS base
type RewriteFS interface {
	RewritePath(string) string
}
