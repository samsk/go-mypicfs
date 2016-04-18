package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	picfs "github.com/samsk/go-mypicfs/fs"
	"github.com/samsk/go-mypicfs/picasa"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

func indentifyFS(sourceDir string) (fsType string) {
	if (picasa.IdentifyPicasaFS(sourceDir)) {
		fsType = "picasa";
	}
	return fsType
}

func main() {
	// Scans the arg list and sets up flags
	var debug = flag.Bool("debug", false, "print debugging messages")
	var fsType = flag.String("type", "", "picfs type (loop, picasa)")
//	other := flag.Bool("allow-other", false, "mount with -o allowother")
	flag.Parse()
	if (flag.NArg() < 2) {
		// TODO - where to get program name?
		fmt.Println("usage: main [args] SOURCE_DIR MOUNT_DIR")
		os.Exit(1)
	}

	sourceDir := flag.Arg(0)
	sourceDir, _ = filepath.Abs(sourceDir)
	mountDir := flag.Arg(1)

	// auto ident FS type
	if (*fsType == "") {
		var ident = indentifyFS(sourceDir)
		if (ident != "") {
			*fsType = ident
		}
	}

	// create object
	var ourFS pathfs.FileSystem
	switch (*fsType) {
		case "loop":
			var fs = picfs.NewLoopFS(sourceDir)
			ourFS = &fs
		case "picasa":
			var fs = picasa.NewPicasaFS(sourceDir)
			ourFS = &fs
		default:
			fmt.Println("ERROR: Filesystem '" + *fsType + "' not supported (yet) !")
			os.Exit(1)
	};

	var opts = &nodefs.Options{
		// These options are to be compatible with libfuse defaults,
		// making benchmarking easier.
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	};
	var pathFs = pathfs.NewPathNodeFs(ourFS, nil)
	var connector = nodefs.NewFileSystemConnector(pathFs.Root(), opts)

	var mountOpts = &fuse.MountOptions{
		AllowOther: true,
		Name:       "mypicfs:" + *fsType,
		FsName:     sourceDir,
	}
	var state, err = fuse.NewServer(connector.RawFS(), mountDir, mountOpts)
	if err != nil {
		fmt.Printf("Mount fail: %v\n", err)
		os.Exit(1)
	}
	state.SetDebug(*debug)

	//fmt.Println("Mounted!")
	state.Serve()
}
