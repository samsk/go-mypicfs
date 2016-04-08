package fs;

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
//	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
);


// Context interface
type Context interface {}

// ContextFS interface
// this is pathfs.FileSystem like interface
type ContextFS interface {
	// Attributes
	CGetAttr(name string, fctx *fuse.Context, ctx Context) (*fuse.Attr, fuse.Status)

	// These should update the file's ctime too.
	CChmod(name string, mode uint32, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CChown(name string, uid uint32, gid uint32, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CUtimens(name string, Atime *time.Time, Mtime *time.Time, fctx *fuse.Context, ctx Context) (code fuse.Status)

	CTruncate(name string, size uint64, fctx *fuse.Context, ctx Context) (code fuse.Status)

	CAccess(name string, mode uint32, fctx *fuse.Context, ctx Context) (code fuse.Status)

	// Tree structure
	CLink(oldName string, newName string, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CMkdir(name string, mode uint32, fctx *fuse.Context, ctx Context) fuse.Status
	CMknod(name string, mode uint32, dev uint32, fctx *fuse.Context, ctx Context) fuse.Status
	CRename(oldName string, newName string, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CRmdir(name string, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CUnlink(name string, fctx *fuse.Context, ctx Context) (code fuse.Status)

	// Extended attributes.
	CGetXAttr(name string, attribute string, fctx *fuse.Context, ctx Context) (data []byte, code fuse.Status)
	CListXAttr(name string, fctx *fuse.Context, ctx Context) (attributes []string, code fuse.Status)
	CRemoveXAttr(name string, attr string, fctx *fuse.Context, ctx Context) fuse.Status
	CSetXAttr(name string, attr string, data []byte, flags int, fctx *fuse.Context, ctx Context) fuse.Status

	// File handling.  If opening for writing, the file's mtime
	// should be updated too.
	COpen(name string, flags uint32, fctx *fuse.Context, ctx Context) (file nodefs.File, code fuse.Status)
	CCreate(name string, flags uint32, mode uint32, fctx *fuse.Context, ctx Context) (file nodefs.File, code fuse.Status)

	// Directory handling
	COpenDir(name string, fctx *fuse.Context, ctx Context) (stream []fuse.DirEntry, code fuse.Status)

	// Symlinks.
	CSymlink(value string, linkName string, fctx *fuse.Context, ctx Context) (code fuse.Status)
	CReadlink(name string, fctx *fuse.Context, ctx Context) (string, fuse.Status)
}
