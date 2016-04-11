# myPicfs

myPicfs is FUSE filesystem wrapper written in go for local Picasa folders. This way it is possible to expose i.e. Picasa starred photos on network via SAMBA or DLNA.

## Mounting

Usage is easy, simply mount directory with your picasa images (i.e. */media/ARCHIVE/photos*) to some mount point (i.e. */media/PHOTOS*).

```bash
$ mypicfs -type=picasa /media/ARCHIVE/photos /media/PHOTOS/
```

## Features

- **Starred photos** per subdirectory
- **Global Starred photos** for all subfolders divided by year

## Planned features

- **Albums** per subdirectory
- **Global Albums** for all subfolders
- On the fly **applying of image filters** (to show them as they were modified in Picasa)

## Dependencies

-  **[gopkg.in/ini.v1](https://gopkg.in/ini.v1)**
-  **[github.com/hanwen/go-fuse](https://github.com/hanwen/go-fuse)**
-  **FUSE** support in kernel

## Building

In go-mypicfs/mypics directory execute

```bash
$ go get
$ go build
```

## Notice

- Only **.picasa.ini** files are recognized and interpreted, you should rename all existing **Picasa.ini** files to make it all work.
- Right now **.picasa.ini** files are cached for 1/2 hour, therefore it might take up to 1/2 hour till your changes (i.e. adding a starred photo) will propagate.
