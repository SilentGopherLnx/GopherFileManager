# What is this?
**GopherFileManager** is file manager for Linux with GTK3 user interface written on Golang

I wanted to create **Linux file manager with "folder trumbnails"** like in Dolpin file manager

Source code is so bad now. It will be added later... Look here now https://github.com/SilentGopherLnx/screenshots_and_binaries/tree/master/BIN64_GopherFileManagerFileMoverGui

![screenshot](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/manager_01.png)

# Dependencies for GOPATH:
1) Golang GTK3 lib
https://github.com/gotk3/gotk3
also for gtk:
> sudo apt-get install libgtk-3-dev
>
> sudo apt-get install libcairo2-dev
>
> sudo apt-get install libglib2.0-dev
2) my Framework
https://github.com/SilentGopherLnx/easygolang

3) **inotify** wrapper for directory update events
- github.com/fsnotify/fsnotify
- https://github.com/golang/sys


4) some libs for additional **image trumbnails** support (webp, ico, bmp ...)
- https://github.com/golang/image
- https://github.com/biessek/golang-ico
- https://github.com/jsummers/gobmp

# Dependencies of programs:
1) Compiled version of my "FileMoverGui" written on golang too
https://github.com/SilentGopherLnx/FileMoverGui

2) **ffmpeg** for **video trumbnails**
> sudo apt-get install ffmpeg

# Status:
App is under development and not versionized

**Not all functions are implemented and realised as planned!**

License type is GPL3

# Platform:
Only Linux. Tested only on amd64 on Cinnamon desktop.

Windows support is NOT planned

# Some configurations for good work:
Open "View/Options" and set:
 - Path to **FileMoverGui** app on golang
 - Path to **hash folder** for trumbnails. Create it! 
 - System file manager name
 - System terminal name
