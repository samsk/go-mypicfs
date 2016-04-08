package picasa

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"path/filepath"

	picfs "devel.dob.sk/go-mypicfs/fs"
	"devel.dob.sk/go-mypicfs/fs/ruleFS"


	"github.com/hanwen/go-fuse/fuse"
//	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
);

// PicasaFS definition
type PicasaFS struct {
	picfs.RuleFS
	ini picasaINI
	cache cacheMap

	loopFS picfs.LoopFS
}

type rootCache struct {
	iniFiles []string
}

func NewPicasaFS(root string) PicasaFS {
	var loopFS = picfs.NewLoopFS(root)
	var rules = ruleFS.NewRegexpRules()

	var fs = PicasaFS {
		RuleFS: picfs.NewRuleFS(&rules, &loopFS),
		ini: newPicasaINI(),
		cache: newCacheMap(),

		loopFS: loopFS,
	}

	// readonly
	fs.SetReadOnly(true)

	// setup rules
	rules.Match(RULE_MATCH_ROOT, RULE_IDENT_ROOT, &fs)
	rules.Match(RULE_MATCH_ROOT_STARRED_DIR, RULE_IDENT_ROOT_STARRED_DIR, &fs)
	rules.Match(RULE_MATCH_ROOT_STARRED_YEAR, RULE_IDENT_ROOT_STARRED_YEAR, &fs)
	rules.Match(RULE_MATCH_ROOT_STARRED_YEAR_FILE, RULE_IDENT_ROOT_STARRED_YEAR_FILE, &fs)
	rules.Match(RULE_MATCH_SUBDIR_STARRED_DIR, RULE_IDENT_SUBDIR_STARRED_DIR, &fs)
	rules.Match(RULE_MATCH_SUBDIR_STARRED_FILE, RULE_IDENT_SUBDIR_STARRED_FILE, &fs)
	// catch all rule
	rules.MatchDefault(-1, &fs)
	return fs
}

// PicasaFS implementation
func (fs *PicasaFS) getVPath(name string, data *ruleFS.RegexpRuleData) (string, uint32) {
	switch (data.Ident) {
		case RULE_IDENT_ROOT_STARRED_DIR:
			fallthrough
		case RULE_IDENT_ROOT_STARRED_YEAR:
			return PICASA_INI_FILE, fuse.S_IFDIR
		case RULE_IDENT_ROOT_STARRED_YEAR_FILE:
			return PICASA_INI_FILE, fuse.S_IFLNK
		case RULE_IDENT_SUBDIR_STARRED_DIR:
			return data.StringSubmatch[1] + "/" + PICASA_INI_FILE, fuse.S_IFDIR
		case RULE_IDENT_SUBDIR_STARRED_FILE:
			return data.StringSubmatch[1] + "/" + PICASA_INI_FILE, fuse.S_IFLNK
	}
//	panic(fmt.Sprintf("getVPath(%s, %p) - unhandled case\n", name, data))
	return "", 0
}

//func (fs *PicasaFS) SetDebug(debug bool) {
//	picfs.RegexpRuleFS.SetDebug(debug)
//}

func (fs *PicasaFS) String() string {
	return fmt.Sprintf("PicasaFS(%s,%s)", fs.FileSystem.String(), fs.loopFS.Root)
}

func (fs *PicasaFS) StatFs(name string) *fuse.StatfsOut {
	return nil
}

func (fs *PicasaFS) CGetAttr(name string, fctx *fuse.Context, ctx picfs.Context) (attr *fuse.Attr, status fuse.Status) {
	var data = ctx.(ruleFS.RegexpRuleData);

	switch (data.Ident) {
		case RULE_IDENT_ROOT_STARRED_DIR:
			fallthrough
		case RULE_IDENT_ROOT_STARRED_YEAR:
			fallthrough
		case RULE_IDENT_ROOT_STARRED_YEAR_FILE:
			fallthrough
		case RULE_IDENT_SUBDIR_STARRED_DIR:
			fallthrough
		case RULE_IDENT_SUBDIR_STARRED_FILE:
			var vpath, vmode = fs.getVPath(name, &data)
			attr, status = fs.loopFS.GetAttr(vpath, fctx)
			if (status == fuse.OK) {
				attr.Mode = (attr.Mode ^ fuse.S_IFREG) | vmode | 0755
			}
		default:
			attr, status = fs.loopFS.GetAttr(name, fctx)
	}
	return attr, status
}

func (fs *PicasaFS) CChmod(name string, mode uint32, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CChown(name string, uid uint32, gid uint32, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CUtimens(name string, Atime *time.Time, Mtime *time.Time, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	var data = ctx.(ruleFS.RegexpRuleData);
	var vpath, _ = fs.getVPath(name, &data)

	return fs.loopFS.Utimens(vpath, Atime, Mtime, fctx)
}

func (fs *PicasaFS) CTruncate(name string, size uint64, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CAccess(name string, mode uint32, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	var data = ctx.(ruleFS.RegexpRuleData);

	var _, vmode = fs.getVPath(name, &data)
	if (vmode != 0) {
		// TODO: filter mode (also in RuleFS?)
		return fuse.OK
	} else {
		return fs.loopFS.Access(name, mode, fctx)
	}
}

func (fs *PicasaFS) CLink(oldName string, newName string, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CMkdir(name string, mode uint32, fctx *fuse.Context, ctx picfs.Context) fuse.Status {
	return fuse.EROFS
}

func (fs *PicasaFS) CMknod(name string, mode uint32, dev uint32, fctx *fuse.Context, ctx picfs.Context) fuse.Status {
	return fuse.EACCES
}

func (fs *PicasaFS) CRename(oldName string, newName string, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CRmdir(name string, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CUnlink(name string, fctx *fuse.Context, ctx picfs.Context) (fuse.Status) {
	return fuse.EROFS
}

func (fs *PicasaFS) CGetXAttr(name string, attribute string, fctx *fuse.Context, ctx picfs.Context) ([]byte, fuse.Status) {
	return nil, fuse.ENODATA
}

func (fs *PicasaFS) CListXAttr(name string, fctx *fuse.Context, ctx picfs.Context) ([]string, fuse.Status) {
	return nil, fuse.ENODATA
}

func (fs *PicasaFS) CRemoveXAttr(name string, attr string, fctx *fuse.Context, ctx picfs.Context) fuse.Status {
	return fuse.EROFS
}

func (fs *PicasaFS) CSetXAttr(name string, attr string, data []byte, flags int, fctx *fuse.Context, ctx picfs.Context) fuse.Status {
	return fuse.EROFS
}

func (fs *PicasaFS) COpen(name string, flags uint32, fctx *fuse.Context, ctx picfs.Context) (nodefs.File, fuse.Status) {
	return fs.loopFS.Open(name, flags, fctx)
}

func (fs *PicasaFS) CCreate(name string, flags uint32, mode uint32, fctx *fuse.Context, ctx picfs.Context) (nodefs.File, fuse.Status) {
	return nil, fuse.EROFS
}

func (fs *PicasaFS) openDir(name string, fctx *fuse.Context, data ruleFS.RegexpRuleData) ([]fuse.DirEntry, fuse.Status) {
	streamOri, status := fs.loopFS.OpenDir(name, fctx)
	if (status != fuse.OK) {
		return streamOri, status
	}

	var havePicasaIni = false
	var haveNoMedia = false

	// look for special files
	var stream []fuse.DirEntry
	for _, element := range streamOri {
		if (element.Name == PICASA_INI_FILE) {
			havePicasaIni = true
		} else if (element.Name == NOMEDIA_FILE) {
			// quick exit
			haveNoMedia = true
			break
		} else {
			stream = append(stream, element)
		}
	}

	// append virtual dirs
	if (haveNoMedia) {
		stream = []fuse.DirEntry {}
	} else if (havePicasaIni) {
		stream = append(stream, fuse.DirEntry {
			Name: PICASA_VDIR_prefix + PICASA_VDIR_STARRED,
			Mode: fuse.S_IFDIR | 0755,
		})
	}
	return stream, status
}

func (fs *PicasaFS) getSubdirPicasaIniList(name string, fctx *fuse.Context) (list []string) {
	var cache = fs.cache.GetEntry("root");
	if (cache != nil) {
		return cache.(*rootCache).iniFiles
	}

	// get all entries
	var files, code = fs.loopFS.OpenDir(name, fctx)
	if (code != fuse.OK || files == nil) {
		return list
	}

	// filter dirs
	for _, element := range files {
		if (element.Mode & fuse.S_IFDIR != 0) {
			// .nomedia
			var file = fs.loopFS.RewritePath(name + "/" + element.Name + "/" + NOMEDIA_FILE)
			if (fs.loopFS.FileSystem.Access(file, fuse.R_OK, fctx) == fuse.OK) {
				continue
			}

			// .picasa.ini
			file = fs.loopFS.RewritePath(name + "/" + element.Name + "/" + PICASA_INI_FILE)
			if (fs.loopFS.FileSystem.Access(file, fuse.R_OK, fctx) == fuse.OK) {
				list = append(list, file)
			}
		}
	}
	fs.cache.SetEntry("root", &rootCache {
		iniFiles: list,
	})
	return list
}

func (fs *PicasaFS) openDir_RootStarred(name string, fctx *fuse.Context, data ruleFS.RegexpRuleData) ([]fuse.DirEntry, fuse.Status) {
	var files = fs.getSubdirPicasaIniList("", fctx)

	var entries []fuse.DirEntry
	var years = make(map[int]bool)
	for _, element := range files {
		var ini = fs.ini.load(element);
		var tim = fs.ini.getDate(ini, "Picasa=>date");

		var year = tim.Year()
		if (years[year] == true) {
			continue
		}
		years[year] = true

		entries = append(entries, fuse.DirEntry {
			Name: fmt.Sprintf("%04d", year),
			Mode: fuse.S_IFDIR | 0755,
		})
	}
	return entries, fuse.OK
}

func (fs *PicasaFS) openDir_RootStarredYear(name string, fctx *fuse.Context, data ruleFS.RegexpRuleData) ([]fuse.DirEntry, fuse.Status) {
	var yearDir, _ = strconv.Atoi(data.StringSubmatch[1])
	var files = fs.getSubdirPicasaIniList("", fctx)

	var entries []fuse.DirEntry
	for _, element := range files {
		var dir = filepath.Base(filepath.Dir(element))
		var ini = fs.ini.load(element);
		var tim = fs.ini.getDate(ini, "Picasa=>date");

		if (tim.Year() != yearDir) {
			continue
		}

		var starred = fs.ini.getStarred(ini);
		for _, file := range starred {
				entries = append(entries, fuse.DirEntry {
				Name: fmt.Sprintf("%s__%s", dir, file),
				Mode: fuse.S_IFLNK | 0755,
			});
		}
	}
	return entries, fuse.OK
}

func (fs *PicasaFS) openDir_SubdirStarred(name string, fctx *fuse.Context, data ruleFS.RegexpRuleData) ([]fuse.DirEntry, fuse.Status) {
	var vpath, _ = fs.getVPath(name, &data)
	var ini = fs.ini.load(fs.loopFS.RewritePath(vpath));

	var stream []fuse.DirEntry
	var starred = fs.ini.getStarred(ini);
	for _, file := range starred {
		stream = append(stream, fuse.DirEntry {
			Name: file,
			Mode: fuse.S_IFLNK | 0755,
		});
	}
	return stream, fuse.OK
}

func (fs *PicasaFS) COpenDir(name string, fctx *fuse.Context, ctx picfs.Context) ([]fuse.DirEntry, fuse.Status) {
	var data = ctx.(ruleFS.RegexpRuleData);

	switch data.Ident {
		case RULE_IDENT_ROOT_STARRED_DIR:
			return fs.openDir_RootStarred(name, fctx, data)
		case RULE_IDENT_ROOT_STARRED_YEAR:
			return fs.openDir_RootStarredYear(name, fctx, data)
		case RULE_IDENT_SUBDIR_STARRED_DIR:
			return fs.openDir_SubdirStarred(name, fctx, data)
		default:
			return fs.openDir(name, fctx, data)
	}
}

func (fs *PicasaFS) CSymlink(value string, linkName string, fctx *fuse.Context, ctx picfs.Context) fuse.Status {
	return fuse.EROFS
}

func (fs *PicasaFS) CReadlink(name string, fctx *fuse.Context, ctx picfs.Context) (string, fuse.Status) {
	var data = ctx.(ruleFS.RegexpRuleData);

	switch (data.Ident) {
		case RULE_IDENT_SUBDIR_STARRED_FILE:
			var path = "../" + data.StringSubmatch[2]
			return path, fuse.OK
		case RULE_IDENT_ROOT_STARRED_YEAR_FILE:
			var file = filepath.Base(name)
			var str = strings.Split(file, "__")
			var path = "../../" + str[0] + "/" + str[1];
			return path, fuse.OK
	}
	return "", fuse.ENOSYS
}
