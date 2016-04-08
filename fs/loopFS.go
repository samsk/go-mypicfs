package fs;
////
// Basic LoopFS with path rewrite capatibility
//

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
);

// FS definition
type LoopFS struct {
	pathfs.FileSystem
//	RewriteFS
	Root string
	ReadOnly bool
}

func NewLoopFS(root string) LoopFS {
	return LoopFS {
		FileSystem: pathfs.NewLoopbackFileSystem("/"),
		Root: root,
		ReadOnly: false,
	}
}

func (fs *LoopFS) RewritePath(path string) string {
	var pathOut = filepath.Join(fs.Root, path);

//	fmt.Println("rewrite(" + path + ") -> " + pathOut);
	return pathOut;
}

// LoopFS implementation
func (fs *LoopFS) SetReadOnly(readOnly bool) {
	fs.ReadOnly = readOnly
}

func (fs *LoopFS) SetDebug(debug bool) {
	fs.FileSystem.SetDebug(debug)
}

func (fs *LoopFS) GetAttr(name string, ctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	return fs.FileSystem.GetAttr(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) Readlink(name string, ctx *fuse.Context) (string, fuse.Status) {
	return fs.FileSystem.Readlink(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) Mknod(name string, mode uint32, dev uint32, ctx *fuse.Context) fuse.Status {
	return fs.FileSystem.Mknod(fs.RewritePath(name), mode, dev, ctx)
}

func (fs *LoopFS) Mkdir(name string, mode uint32, ctx *fuse.Context) fuse.Status {
	return fs.FileSystem.Mkdir(fs.RewritePath(name), mode, ctx)
}

func (fs *LoopFS) Unlink(name string, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Unlink(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) Rmdir(name string, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Rmdir(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) Symlink(value string, linkName string, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Symlink(value, fs.RewritePath(linkName), ctx)
}

func (fs *LoopFS) Rename(oldName string, newName string, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Rename(fs.RewritePath(oldName), fs.RewritePath(newName), ctx)
}

func (fs *LoopFS) Link(oldName string, newName string, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Link(fs.RewritePath(oldName), fs.RewritePath(newName), ctx)
}

func (fs *LoopFS) Chmod(name string, mode uint32, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Chmod(fs.RewritePath(name), mode, ctx)
}

func (fs *LoopFS) Chown(name string, uid uint32, gid uint32, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Chown(fs.RewritePath(name), uid, gid, ctx)
}

func (fs *LoopFS) Truncate(name string, offset uint64, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Truncate(fs.RewritePath(name), offset, ctx)
}

func (fs *LoopFS) Open(name string, flags uint32, ctx *fuse.Context) (file nodefs.File, code fuse.Status) {
	return fs.FileSystem.Open(fs.RewritePath(name), flags, ctx)
}

func (fs *LoopFS) OpenDir(name string, ctx *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	return fs.FileSystem.OpenDir(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) OnMount(nodeFs *pathfs.PathNodeFs) {
	fs.FileSystem.OnMount(nodeFs)
}

func (fs *LoopFS) OnUnmount() {
	fs.FileSystem.OnUnmount()
}

func (fs *LoopFS) Access(name string, mode uint32, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Access(fs.RewritePath(name), mode, ctx)
}

func (fs *LoopFS) Create(name string, flags uint32, mode uint32, ctx *fuse.Context) (file nodefs.File, code fuse.Status) {
	return fs.FileSystem.Create(fs.RewritePath(name), flags, mode, ctx)
}

func (fs *LoopFS) Utimens(name string, Atime *time.Time, Mtime *time.Time, ctx *fuse.Context) (code fuse.Status) {
	return fs.FileSystem.Utimens(fs.RewritePath(name), Atime, Mtime, ctx)
}

func (fs *LoopFS) GetXAttr(name string, attr string, ctx *fuse.Context) ([]byte, fuse.Status) {
	return fs.FileSystem.GetXAttr(fs.RewritePath(name), attr, ctx)
}

func (fs *LoopFS) SetXAttr(name string, attr string, data []byte, flags int, ctx *fuse.Context) fuse.Status {
	return fs.FileSystem.SetXAttr(fs.RewritePath(name), attr, data, flags, ctx)
}

func (fs *LoopFS) ListXAttr(name string, ctx *fuse.Context) ([]string, fuse.Status) {
	return fs.FileSystem.ListXAttr(fs.RewritePath(name), ctx)
}

func (fs *LoopFS) RemoveXAttr(name string, attr string, ctx *fuse.Context) fuse.Status {
	return fs.FileSystem.RemoveXAttr(fs.RewritePath(name), attr, ctx)
}

func (fs *LoopFS) String() string {
	return fmt.Sprintf("LoopFS(%s,%s)", fs.FileSystem.String(), fs.Root)
}

func (fs *LoopFS) StatFs(name string) *fuse.StatfsOut {
	return fs.FileSystem.StatFs(fs.RewritePath(name))
}
