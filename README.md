# What is this?
**GopherFileManager** is file manager for Linux with GTK3 user interface written on Golang

I wanted to create **Linux file manager with "folder trumbnails"** like in Dolphin file manager

Source code is so bad now. It maybe will fix it later... Look here for executable file https://github.com/SilentGopherLnx/screenshots_and_binaries/tree/master/BIN64_GopherFileManagerFileMoverGui

![screenshot](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/manager_00.png)

![screenshot](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/manager_01.png)

![screenshot](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/manager_02.png)

![screenshot](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/manager_03.png)

more screenshoots:
> https://github.com/SilentGopherLnx/screenshots_and_binaries/tree/master/SCREENS_GopherFileManagerFileMoverGui

# Dependencies for GOPATH (can be downloaded by "go get"):
1) Golang **GTK3** lib 0.6.1
https://github.com/gotk3/gotk3
also for gtk:
> sudo apt-get install libgtk-3-dev
>
> sudo apt-get install libcairo2-dev
>
> sudo apt-get install libglib2.0-dev
2) **my Framework**
https://github.com/SilentGopherLnx/easygolang

3) **inotify** wrapper for directory update events
- https://github.com/fsnotify/fsnotify
- https://github.com/golang/sys

4) some libs for additional **image trumbnails** support (webp, ico, bmp ...)
- https://github.com/golang/image
- https://github.com/biessek/golang-ico
- https://github.com/jsummers/gobmp

5) correct image rotation from exif
- https://github.com/disintegration/imageorient
- https://github.com/disintegration/gift

# Dependencies of PROGRAMS (download by yourself!):
1) Compiled version of my **"FileMoverGui"** written on golang too
https://github.com/SilentGopherLnx/FileMoverGui

2) **ffmpeg** for **video trumbnails**
> sudo apt-get install ffmpeg

3) **smbclient** for getting smb folders list of one server
> sudo apt-get install smbclient

**programs below usualy included in linux** (if you not have them, you should install it):
1) gtk, gvfs, bash - libs
2) df, lsblk, mount, ls - base commands for disks list
3) xdg-open, xdg-mime - file type system associations
4) **xclip** - for copy-paste clipboard
5) avahi-browse - for scanning network for samba servers
6) udevadm - for usb devices names (mtp protocol)

# Status:
App is under development and looks freaky

**Not all functions are implemented and realised as planned!** need to do:
1) "hard-to-do":
- ~~**async file list loading** by os.File.Readdir(1+) (ioutil.ReadDir() or os.File.Readdir(-1) is slow if mpt/webdav protocol)~~
- ~~**listing all remote pc on network, list folders of one remote pc**~~
- ~~**searching files** results~~
- ~~mounting remote folders~~ & mounting remote folders with unsaved password & **mounting local (unmounted or encrypted) folders**
2) easy (so, it will be done later):
- ~~history of location ("back button" fo path)~~
- some small features like ~~sorting~~, favorite folders editing
- hash folder automatic clear for too old preview-images
- list style show of file list
- and so on

**possible futures in far-far future:**
- working inside zip files like in directories
- preview for zip files content like in folders
- trash for deleted files
- 4K display support

**NEW:**
- 0.2.0 - added multiple select files by mouse + hotkeys (Ctrl+C,Ctrl+V,Del)
- 0.2.1 - hotkeys F2,F5,Ctrl+A; info for multiseleced files; copy folder near itself
- 0.2.2 - new cache method, exif orientation, paste into in menu, russian keyboard shortcuts fix
- 0.3.0 - async file list loading! searching files! network listing! movie files length on icon
- 0.3.1 - smb mount, path history, search results menus
- 0.3.2 - mouse selecting with Ctrl button, fixing unnecessary menu's bugs
- 0.3.3 - fixing for install and compile from source code by **go get**
- **0.4.0** - multi-language support(+russian), checking for new verison button, fix for latest gtk release, fixing for multi-language keyboard hotkeys (Ctrl+C,Ctrl+X,Ctrl+V),renaming gialog moved to helper app
- **0.4.1** - fix for **gotk3 0.6.1** (was uncompileable in selectionData) and **golang 1.17** (url parsing works differently); add: count of files on folder preview, sort files by size, option to select preview update time (always-never) and app for rename dialog
- **0.4.2** - easier smb mounting, movie height on preview, sorting by file date, folder preview error icon, more language translation

# Platform & License:
**Only Linux!** Tested mostly on amd64 on Cinnamon desktop (Linux Mint ~~19~~20.2). 
GTK 3.24.20
Also tested once on Gentoo 17 (gnome 3) amd64.

Windows/MasOS support is NOT planned

**License type is GPL3**

# Some configurations for good work:
Open "View/Options" and change if you need:
 - Path to **FileMoverGui** app on golang
 - Path to **cache folder** for trumbnails. Create it if not exist or choose another! 
 - System file manager name ("nemo" is default)
 - System terminal name ("gnome-terminal" is default)

