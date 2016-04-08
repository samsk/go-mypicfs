package picasa;
// picasa ini manager

import (
	"fmt"
	"time"
	"strings"
//	"strconv"

	"gopkg.in/ini.v1"
);

type picasaINI struct {
	cache cacheMap
};

type picasaINIcache struct {
	file *ini.File
}

func newPicasaINI() picasaINI {
	return picasaINI {
		cache: newCacheMap(),
	}
}

func (obj *picasaINI) load(path string) (file *ini.File) {
	// try cache first
	// not using stat/mtime to avoid I/O
	var data = obj.cache.GetEntry(path)
	if (data != nil) {
		return data.(*picasaINIcache).file
	} else {
		var err error
		file, err = ini.Load(path);
		if (err != nil) {
			fmt.Printf("ini: file '%s' empty !\n", path)
			file = ini.Empty();
		}

		// tune
		file.BlockMode = false

		// refresh cache
		obj.cache.SetEntry(path, &picasaINIcache {
			file: file,
		})
	}
	return file;
}

func (obj *picasaINI) getStarred(file *ini.File) (files []string) {
	var sections = file.SectionStrings()

	for _, element := range sections {
		var starred = file.Section(element).Key("star").MustBool(false)
		var category = file.Section(element).Key("category").MustString("")

		if (category == "Folders on Disk") {
			continue
		} else if (starred == true) {
			files = append(files, element)
		} 
	}
	return files;
}

func (obj *picasaINI) getKey(file *ini.File, key string) *ini.Key {
	var str = strings.Split(key, "=>")
	var section = ""

	if (len(str) > 1) {
		section = str[0]
		key = str[1]
	}
	return file.Section(section).Key(key)
}

func (obj *picasaINI) getDate(file *ini.File, key string) time.Time {
	var val, _ = obj.getKey(file, key).Float64()

	var dur, _ = time.ParseDuration(fmt.Sprintf("%.6fh", val * 24))
	return time.Date(1900, 0, 0, 0, 0, 0, 0, time.Local).Add(dur);
}
