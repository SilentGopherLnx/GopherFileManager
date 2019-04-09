package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"github.com/gotk3/gotk3/gtk"
)

func Menu_CurrentFolder(menu *gtk.Menu, folderpath string) {
	paste_list := LinuxClipBoard_PasteFiles()
	var func_paste func() = nil
	if len(paste_list) > 0 {
		func_paste = func() {
			GTK_CopyPasteDnd_Paste(folderpath)
		}
	}
	GTK_MenuItem(menu, "Paste (Ctrl+V)", func_paste)
	submenu_new := GTK_MenuSub(menu, "New")
	GTK_MenuItem(submenu_new, "Folder", func() {
		create_new(folderpath, true)
		listFiles(gGFiles, folderpath)
	})
	GTK_MenuSeparator(submenu_new)
	GTK_MenuItem(submenu_new, "Text File", nil)
	GTK_MenuItem(submenu_new, "Empty File", func() {
		create_new(folderpath, false)
		listFiles(gGFiles, folderpath)
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
	GTK_MenuItem(submenu_sort, "INC", nil)
	GTK_MenuItem(submenu_sort, "DECR", nil)
	GTK_MenuSeparator(submenu_sort)
	GTK_MenuItem(submenu_sort, "Name", nil)
	GTK_MenuItem(submenu_sort, "Type", nil)
	GTK_MenuItem(submenu_sort, "Size", nil)
	GTK_MenuSeparator(menu)
	GTK_MenuItem(menu, "Info", nil)
}

func Menu_FilesContextMenu(menu *gtk.Menu, fpath string, fname string, isdir bool, isapp bool) {
	fpath2 := FolderPathEndSlash(fpath)
	ext_mime := FileMIME(fpath2 + fname)
	app_mime := AppMIME(ext_mime)
	apps_mime := AllAppsMIME(ext_mime)
	GTK_MenuItem(rightmenu, Select_String(!isapp, "*", "")+"Open ["+app_mime+"]", nil)
	if isdir {
		GTK_MenuItem(rightmenu, "Open in new Window", func() {
			menu.Cancel()
			go ExecCommandBash("" + ExecQuote(AppRunArgs()[0]) + " " + ExecQuote(fpath2+fname) + "")
		})
		//if(islink){
		GTK_MenuItem(rightmenu, "Open with eval symlink", nil)
	} else {
		GTK_MenuItem(rightmenu, Select_String(isapp, "*", "")+"Run", nil)
		//GTK_MenuSeparator(rightmenu)

		submenu_openwith := GTK_MenuSub(rightmenu, "Open With")
		GTK_MenuItem(submenu_openwith, "Text Editor", func() {
			OpenFileByApp(fpath2+fname, opt.GetTextEditor())
		})
		GTK_MenuItem(submenu_openwith, "HEX Editor", nil)
		GTK_MenuItem(submenu_openwith, "Archive", nil)
		GTK_MenuItem(submenu_openwith, "Copy image to clipboard", nil)
		GTK_MenuItem(submenu_openwith, "Other...", func() {
			//gtk.AppChooserDialogNewForContentType()
		})
		GTK_MenuSeparator(submenu_openwith)
		for j := 0; j < len(apps_mime); j++ {
			app_alt := apps_mime[j]
			GTK_MenuItem(submenu_openwith, app_alt, func() {
				OpenFileByApp(fpath2+fname, app_alt)
			})
		}
		GTK_MenuSeparator(submenu_openwith)
		GTK_MenuItem(submenu_openwith, "?", nil)

		submenu_runas := GTK_MenuSub(rightmenu, "Run AS")
		GTK_MenuItem(submenu_runas, "SUDO", nil)
		GTK_MenuItem(submenu_runas, "in Terminal", nil)
		GTK_MenuItem(submenu_runas, "SUDO in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "BASH in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "SUDO BASH in Terminal", nil)
	}

	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, "Cut (Ctrl+X)", func() {
		copied := NewLinuxPath(false) //??
		copied.SetReal(fpath2 + fname)
		Prln("cut: " + copied.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{copied}, true)
	})
	GTK_MenuItem(rightmenu, "Copy (Ctrl+C)", func() {
		copied := NewLinuxPath(false) //??
		copied.SetReal(fpath2 + fname)
		Prln("copy: " + copied.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{copied}, false)
	})
	var del_func func() = nil
	//if !isdir {
	del_func = func() {
		Dialog_FileDelete(win, fpath, fname, func() {
			listFiles(gGFiles, fpath2)
		})
	}
	//}

	//item :=

	GTK_MenuItemIcon(rightmenu, "Delete (Del)", "delete", del_func)

	//item = item

	// theme, _ := gtk.IconThemeGetDefault()
	// pixbuf, _ := theme.LoadIcon("edit-copy", 16, 0)
	// img, _ := gtk.ImageNewFromPixbuf(pixbuf)
	// wi, _ := item.GetChild()
	// w
	// //lbl := wi.(gtk.accel )

	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, "Rename (F2)", func() {
		Dialog_FileRename(win, fpath2, fname, func() {
			listFiles(gGFiles, fpath2)
		})
	})
	GTK_MenuItem(rightmenu, "Compress", nil)

	GTK_MenuSeparator(rightmenu)
	if isapp {
		GTK_MenuItem(rightmenu, "Create Shortcut", nil)
	}
	if isdir {
		GTK_MenuItem(rightmenu, "Add To Favorites", nil)
	}
	GTK_MenuItem(rightmenu, "Info ["+ext_mime+"]", func() {
		Dialog_FileInfo(win, fpath2, fname)
	})
}

func GTK_MainMenu(win *gtk.Window) *gtk.MenuBar {
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
	Menu_CurrentFolder(submenu_edit, path.GetReal())

	submenu_other := GTK_MenuSub(menuBar, "Info")
	GTK_MenuItem(submenu_other, "Help", nil)
	GTK_MenuItem(submenu_other, "About", func() {
		Dialog_About(win, AppVersion(), AppAuthor(), AppMail(), AppRepository(), GetFlag_Russian())
	})
	return menuBar
}
