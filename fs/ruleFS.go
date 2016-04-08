package fs
////
// Multiple filesystem layers activated by file based rules
//

import (
//	"fmt"
//	"path/filepath"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
);

// rule match
type RuleFSMatch interface {
	// match rule
	MatchPath(name string) (ContextFS, Context)
}

// rule FS definition
type RuleFS struct {
	pathfs.FileSystem
	RuleFSMatch
	ReadOnly bool
}

func NewRuleFS(rule RuleFSMatch, fs pathfs.FileSystem) RuleFS {
	return RuleFS {
		RuleFSMatch: rule,
		FileSystem: fs,
		ReadOnly: false,
	}
}

// FS implementation
func (fs *RuleFS) SetReadOnly(readOnly bool) {
	fs.ReadOnly = readOnly
}

func (fs *RuleFS) GetAttr(name string, fctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CGetAttr(name, fctx, ctx)
	} else {
		return fs.FileSystem.GetAttr(name, fctx)
	}
}

func (fs *RuleFS) Readlink(name string, fctx *fuse.Context) (string, fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CReadlink(name, fctx, ctx)
	} else {
		return fs.FileSystem.Readlink(name, fctx)
	}
}

func (fs *RuleFS) Mknod(name string, mode uint32, dev uint32, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CMknod(name, mode, dev, fctx, ctx)
	} else {
		return fs.FileSystem.Mknod(name, mode, dev, fctx)
	}
}

func (fs *RuleFS) Mkdir(name string, mode uint32, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CMkdir(name, mode, fctx, ctx)
	} else {
		return fs.FileSystem.Mkdir(name, mode, fctx)
	}
}

func (fs *RuleFS) Unlink(name string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CUnlink(name, fctx, ctx)
	} else {
		return fs.FileSystem.Unlink(name, fctx)
	}
}

func (fs *RuleFS) Rmdir(name string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)	
	if (outfs != nil) {
		return outfs.CRmdir(name, fctx, ctx)
	} else {
		return fs.FileSystem.Rmdir(name, fctx)
	}
}

func (fs *RuleFS) Symlink(value string, linkName string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	//XXX: linkName or value ???
	var outfs, ctx = fs.MatchPath(linkName)
	
	if (outfs != nil) {
		return outfs.CSymlink(value, linkName, fctx, ctx)
	} else {
		return fs.FileSystem.Symlink(value, linkName, fctx)
	}
}

func (fs *RuleFS) Rename(oldName string, newName string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(oldName)	
	if (outfs != nil) {
		return outfs.CRename(oldName, newName, fctx, ctx)
	} else {
		return fs.FileSystem.Rename(oldName, newName, fctx)
	}
}

func (fs *RuleFS) Link(oldName string, newName string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(oldName)	
	if (outfs != nil) {
		return outfs.CLink(oldName, newName, fctx, ctx)
	} else {
		return fs.FileSystem.Link(oldName, newName, fctx)
	}
}

func (fs *RuleFS) Chmod(name string, mode uint32, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}
	
	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CChmod(name, mode, fctx, ctx)
	} else {
		return fs.FileSystem.Chmod(name, mode, fctx)
	}
}

func (fs *RuleFS) Chown(name string, uid uint32, gid uint32, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CChown(name, uid, gid, fctx, ctx)
	} else {
		return fs.FileSystem.Chown(name, uid, gid, fctx)
	}
}

func (fs *RuleFS) Truncate(name string, offset uint64, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CTruncate(name, offset, fctx, ctx)
	} else {
		return fs.FileSystem.Truncate(name, offset, fctx)
	}
}

func (fs *RuleFS) Open(name string, flags uint32, fctx *fuse.Context) (file nodefs.File, code fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.COpen(name, flags, fctx, ctx)
	} else {
		return fs.FileSystem.Open(name, flags, fctx)
	}
}

func (fs *RuleFS) OpenDir(name string, fctx *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.COpenDir(name, fctx, ctx)
	} else {
		return fs.FileSystem.OpenDir(name, fctx)
	}
}

func (fs *RuleFS) Access(name string, mode uint32, fctx *fuse.Context) (code fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CAccess(name, mode, fctx, ctx)
	} else {
		return fs.FileSystem.Access(name, mode, fctx)
	}
}

func (fs *RuleFS) Create(name string, flags uint32, mode uint32, fctx *fuse.Context) (file nodefs.File, code fuse.Status) {
	if (fs.ReadOnly) {
		return nil, fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CCreate(name, flags, mode, fctx, ctx)
	} else {
		return fs.FileSystem.Create(name, flags, mode, fctx)
	}
}

func (fs *RuleFS) Utimens(name string, Atime *time.Time, Mtime *time.Time, fctx *fuse.Context) fuse.Status {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CUtimens(name, Atime, Mtime, fctx, ctx)
	} else {
		return fs.FileSystem.Utimens(name, Atime, Mtime, fctx)
	}
}

func (fs *RuleFS) GetXAttr(name string, attr string, fctx *fuse.Context) ([]byte, fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CGetXAttr(name, attr, fctx, ctx)
	} else {
		return fs.FileSystem.GetXAttr(name, attr, fctx)
	}
}

func (fs *RuleFS) SetXAttr(name string, attr string, data []byte, flags int, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CSetXAttr(name, attr, data, flags, fctx, ctx)
	} else {
		return fs.FileSystem.SetXAttr(name, attr, data, flags, fctx)
	}
}

func (fs *RuleFS) ListXAttr(name string, fctx *fuse.Context) ([]string, fuse.Status) {
	var outfs, ctx = fs.MatchPath(name)

	if (outfs != nil) {
		return outfs.CListXAttr(name, fctx, ctx)
	} else {
		return fs.FileSystem.ListXAttr(name, fctx)
	}
}

func (fs *RuleFS) RemoveXAttr(name string, attr string, fctx *fuse.Context) fuse.Status {
	if (fs.ReadOnly) {
		return fuse.EROFS
	}

	var outfs, ctx = fs.MatchPath(name)
	if (outfs != nil) {
		return outfs.CRemoveXAttr(name, attr, fctx, ctx)
	} else {
		return fs.FileSystem.RemoveXAttr(name, attr, fctx)
	}
}
