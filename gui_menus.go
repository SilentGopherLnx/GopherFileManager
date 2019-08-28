package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "./pkg_filetools"

	"github.com/gotk3/gotk3/gtk"
)

func GTKMenu_CurrentFolder(menu *gtk.Menu, folderpath string) {
	paste_list, _ := LinuxClipBoard_PasteFiles()
	var func_paste func() = nil
	if len(paste_list) > 0 {
		func_paste = func() {
			GTK_CopyPasteDnd_Paste(folderpath)
		}
	}
	GTK_MenuItem(menu, "Paste "+I2S(len(paste_list))+"objects (Ctrl+V)", func_paste)
	submenu_new := GTK_MenuSub(menu, "New")
	GTK_MenuItem(submenu_new, "Folder", func() {
		name_created := FileOrFolder_New(folderpath, true)
		listFiles(gGFiles, folderpath, false)
		Dialog_FileRename(win, folderpath, name_created, func() {
			listFiles(gGFiles, folderpath, false)
		})
	})
	GTK_MenuSeparator(submenu_new)
	GTK_MenuItem(submenu_new, "Text File", nil)
	GTK_MenuItem(submenu_new, "Empty File", func() {
		name_created := FileOrFolder_New(folderpath, false)
		listFiles(gGFiles, folderpath, false)
		Dialog_FileRename(win, folderpath, name_created, func() {
			listFiles(gGFiles, folderpath, false)
		})
	})
	GTK_MenuSeparator(menu)
	submenu_refolder := GTK_MenuSub(menu, "Reopen folder")
	GTK_MenuItem(submenu_refolder, "SUDO", nil)
	GTK_MenuItem(submenu_refolder, "in Terminal", func() {
		menu.Cancel()
		term := opt.GetTerminal(folderpath)
		go ExecCommandBash(term)
		//a, b, c :=   Prln(a + b + c)
	})
	GTK_MenuItem(submenu_refolder, "SUDO in Terminal", nil)
	GTK_MenuItem(submenu_refolder, "Default File Manager", func() {
		menu.Cancel()
		fm := opt.GetFileManager(folderpath)
		go ExecCommandBash(fm)
		//a, b, c :=   Prln(a + b + c)
	})
	submenu_sort := GTK_MenuSub(menu, "Sort")
	GTK_MenuItem(submenu_sort, "INC "+B2S(sort_reverse, "", "(v)"), func() {
		sort_reverse = false
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, "DECR "+B2S(sort_reverse, "(v)", ""), func() {
		sort_reverse = true
		resort_and_show()
	})
	GTK_MenuSeparator(submenu_sort)
	GTK_MenuItem(submenu_sort, "Name "+B2S(sort_mode == 0, "(v)", ""), func() {
		sort_mode = 0
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, "Type "+B2S(sort_mode == 1, "(v)", ""), func() {
		sort_mode = 1
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, "Size "+B2S(sort_mode == 2, "(v)", ""), nil)
	GTK_MenuSeparator(menu)
	GTK_MenuItem(menu, "Info", func() {
		Dialog_FileInfo(win, LinuxFileGetParent(folderpath), []string{FolderPathEndSlash(LinuxFileNameFromPath(folderpath))})
	})
}

func GTKMenu_File(menu *gtk.Menu, fpath string, fname string, isdir bool, isapp bool) {
	fpath2 := FolderPathEndSlash(fpath)
	ext_mime := FileMIME(fpath2 + fname)
	app_mime := AppMIME(ext_mime)
	apps_mime := AllAppsMIME(ext_mime)
	GTK_MenuItem(rightmenu, B2S(!isapp, "*", "")+"Open ["+app_mime+"]", nil)
	if isdir {
		GTK_MenuItem(rightmenu, "Open in new Window", func() {
			menu.Cancel()
			go ExecCommandBash("" + ExecQuote(AppRunArgs()[0]) + " " + ExecQuote(fpath2+fname) + B2S(win.IsMaximized(), " -max", ""))
		})
		//if(islink){
		GTK_MenuItem(rightmenu, "Open with eval symlink", nil)
	} else {
		GTK_MenuItem(rightmenu, B2S(isapp, "*", "")+"Run", nil)
		//GTK_MenuSeparator(rightmenu)

		submenu_openwith := GTK_MenuSub(rightmenu, "Open With")
		GTK_MenuItem(submenu_openwith, "Text Editor", func() {
			OpenFileByApp(fpath2+fname, opt.GetTextEditor())
		})
		GTK_MenuItem(submenu_openwith, "HEX Editor", nil)
		GTK_MenuItem(submenu_openwith, "Archive", nil)
		GTK_MenuItem(submenu_openwith, "Copy image to clipboard", nil)
		GTK_MenuSeparator(submenu_openwith)
		for j := 0; j < len(apps_mime); j++ {
			app_alt := apps_mime[j]
			GTK_MenuItem(submenu_openwith, app_alt, func() {
				OpenFileByApp(fpath2+fname, app_alt)
			})
		}
		GTK_MenuSeparator(submenu_openwith)
		GTK_MenuItem(submenu_openwith, "Other?...", nil) // //gtk.AppChooserDialogNewForContentType()

		submenu_runas := GTK_MenuSub(rightmenu, "Run AS")
		GTK_MenuItem(submenu_runas, "SUDO", nil)
		GTK_MenuItem(submenu_runas, "in Terminal", nil)
		GTK_MenuItem(submenu_runas, "SUDO in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "BASH in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "SUDO BASH in Terminal", nil)
	}

	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, "Cut (Ctrl+X)", func() {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fname)
		Prln("cut: " + file1.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{file1}, true)
	})
	GTK_MenuItem(rightmenu, "Copy (Ctrl+C)", func() {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fname)
		Prln("copy: " + file1.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{file1}, false)
	})
	GTK_MenuItemIcon(rightmenu, "Delete (Del)", "delete", func() {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fname)
		Prln("del: " + file1.GetUrl())
		RunFileOperaion([]*LinuxPath{file1}, nil, OPER_DELETE)
	})

	//item = item

	// theme, _ := gtk.IconThemeGetDefault()
	// pixbuf, _ := theme.LoadIcon("edit-copy", 16, 0)
	// img, _ := gtk.ImageNewFromPixbuf(pixbuf)
	// wi, _ := item.GetChild()
	// w
	// //lbl := wi.(gtk.accel )

	GTK_MenuSeparator(rightmenu)
	if isdir {
		paste_list, _ := LinuxClipBoard_PasteFiles()
		var func_paste func() = nil
		if len(paste_list) > 0 {
			func_paste = func() {
				//Prln("[" + fpath2 + fname + "]")
				GTK_CopyPasteDnd_Paste(FolderPathEndSlash(fpath2 + fname))
			}
		}
		GTK_MenuItem(rightmenu, "Paste INTO "+I2S(len(paste_list))+"objects", func_paste) // (Ctrl+V)
	}
	GTK_MenuItem(rightmenu, "Rename (F2)", func() {
		Dialog_FileRename(win, fpath2, fname, func() {
			listFiles(gGFiles, fpath2, false)
		})
	})
	if isdir {
		GTK_MenuItem(rightmenu, "Clear", func() {
			file1 := NewLinuxPath(false) //??
			file1.SetReal(fpath2 + fname)
			Prln("del: " + file1.GetUrl())
			RunFileOperaion([]*LinuxPath{file1}, nil, OPER_CLEAR)
		})
	}
	GTK_MenuItem(rightmenu, "Compress", nil)
	GTK_MenuSeparator(rightmenu)
	if isapp {
		GTK_MenuItem(rightmenu, "Create Shortcut", nil)
	}
	if isdir {
		GTK_MenuItem(rightmenu, "Add To Favorites", nil)
		GTK_MenuItem(rightmenu, "Clear inside", nil)
	}
	GTK_MenuItem(rightmenu, "Info ["+ext_mime+"]", func() {
		fname2 := fname
		if isdir {
			fname2 = FolderPathEndSlash(fname)
		}
		Dialog_FileInfo(win, fpath2, []string{fname2})
	})
}

func GTKMenu_Files(menu *gtk.Menu, fpath string, fnames []string, isdir bool, isapp bool) {
	fpath2 := FolderPathEndSlash(fpath)
	list := []*LinuxPath{}
	for j := 0; j < len(fnames); j++ {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fnames[j])
		list = append(list, file1)
	}
	GTK_MenuItem(rightmenu, "Cut (Ctrl+X)", func() {
		Prln("cut: " + I2S(len(fnames)) + "files")
		LinuxClipBoard_CopyFiles(list, true)
	})
	GTK_MenuItem(rightmenu, "Copy "+I2S(len(fnames))+"objects (Ctrl+C)", func() {
		Prln("copy: " + I2S(len(fnames)) + "files")
		LinuxClipBoard_CopyFiles(list, false)
	})
	GTK_MenuItemIcon(rightmenu, "Delete (Del)", "delete", func() {
		RunFileOperaion(list, nil, OPER_DELETE)
	})
	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, "Info", func() {
		Dialog_FileInfo(win, fpath2, fnames)
	})
}

func GTKMenu_Main(win *gtk.Window) *gtk.MenuBar {
	menuBar, _ := gtk.MenuBarNew()
	submenu_file := GTK_MenuSub(menuBar, "Commands")
	GTK_MenuItem(submenu_file, "New window", nil)
	GTK_MenuItem(submenu_file, "Mount remote folder", nil)
	GTK_MenuItem(submenu_file, "Reload drives list", func() {
		listDiscs(gGDiscs)
	})
	GTK_MenuItem(submenu_file, "Search", nil)
	submenu_view := GTK_MenuSub(menuBar, "View")
	GTK_MenuItem(submenu_view, "Hidden files", nil)
	GTK_MenuItem(submenu_view, "List of files/Icons table", nil)
	GTK_MenuSeparator(submenu_view)
	GTK_MenuItem(submenu_view, "Options", func() {
		Dialog_Options(win)
	})

	submenu_edit := GTK_MenuSub(menuBar, "Current Folder")
	GTKMenu_CurrentFolder(submenu_edit, path.GetReal())

	submenu_other := GTK_MenuSub(menuBar, "Info")
	GTK_MenuItem(submenu_other, "Help", nil)
	GTK_MenuItem(submenu_other, "About", func() {
		Dialog_About(win, AppVersion(), AppAuthor(), AppMail(), AppRepository(), GetFlag_Russian())
	})
	return menuBar
}
