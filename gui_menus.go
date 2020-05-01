package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_filetools"

	"github.com/gotk3/gotk3/gtk"
)

func GTKMenu_CurrentFolder(menu *gtk.Menu, folderpath LinuxPath, search string) {
	GTK_Childs(menu, true, true)
	url := folderpath.GetUrl()
	is_smb, pc_name, netfolder := SMB_CheckVirtualPath(url)
	if StringLength(search) > 0 || is_smb || (StringLength(pc_name) > 0 && StringLength(netfolder) == 0) {
		return
	}
	paste_list, _ := LinuxClipBoard_PasteFiles()
	fp := folderpath.GetReal()
	var func_paste func() = nil
	if len(paste_list) > 0 {
		func_paste = func() {
			GTK_CopyPasteDnd_Paste(fp)
		}
	}
	GTK_MenuItem(menu, langs.GetStr("cmd_folder_paste")+" "+I2S(len(paste_list))+langs.GetStr("cmd_of_objects")+" (Ctrl+V)", func_paste)
	submenu_new := GTK_MenuSub(menu, langs.GetStr("cmd_create"))
	GTK_MenuItem(submenu_new, langs.GetStr("cmd_create_folder"), func() {
		name_created := FileOrFolder_New(fp, true)
		//listFiles(gGFiles, &folderpath, false)
		Dialog_FileRename(win, fp, name_created, func() {
			//listFiles(gGFiles, &folderpath, false)
		})
	})
	GTK_MenuSeparator(submenu_new)
	//GTK_MenuItem(submenu_new, "Text File", nil)
	GTK_MenuItem(submenu_new, langs.GetStr("cmd_create_emptyfile"), func() {
		name_created := FileOrFolder_New(fp, false)
		//listFiles(gGFiles, &folderpath, false)
		Dialog_FileRename(win, fp, name_created, func() {
			//listFiles(gGFiles, &folderpath, false)
		})
	})
	GTK_MenuSeparator(menu)
	submenu_refolder := GTK_MenuSub(menu, "Reopen folder")
	GTK_MenuItem(submenu_refolder, "SUDO", nil)
	GTK_MenuItem(submenu_refolder, "in Terminal", func() {
		menu.Cancel()
		term := opt.GetTerminal(fp)
		go ExecCommandBash(term)
		//a, b, c :=   Prln(a + b + c)
	})
	GTK_MenuItem(submenu_refolder, "SUDO in Terminal", nil)
	GTK_MenuItem(submenu_refolder, "Default File Manager", func() {
		menu.Cancel()
		fm := opt.GetFileManager(fp)
		go ExecCommandBash(fm)
		//a, b, c :=   Prln(a + b + c)
	})
	submenu_sort := GTK_MenuSub(menu, langs.GetStr("cmd_sort"))
	GTK_MenuItem(submenu_sort, langs.GetStr("cmd_sort_inc")+" "+B2S(sort_reverse, "", "(v)"), func() {
		sort_reverse = false
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, langs.GetStr("cmd_sort_decr")+" "+B2S(sort_reverse, "(v)", ""), func() {
		sort_reverse = true
		resort_and_show()
	})
	GTK_MenuSeparator(submenu_sort)
	GTK_MenuItem(submenu_sort, langs.GetStr("cmd_sort_name")+" "+B2S(sort_mode == 0, "(v)", ""), func() {
		sort_mode = 0
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, langs.GetStr("cmd_sort_type")+" "+B2S(sort_mode == 1, "(v)", ""), func() {
		sort_mode = 1
		resort_and_show()
	})
	GTK_MenuItem(submenu_sort, langs.GetStr("cmd_sort_size")+" "+B2S(sort_mode == 2, "(v)", ""), nil)
	GTK_MenuSeparator(menu)
	GTK_MenuItem(menu, langs.GetStr("cmd_files_info"), func() {
		Dialog_FileInfo(win, LinuxFileGetParent(fp), []string{FolderPathEndSlash(LinuxFileNameFromPath(fp))})
	})
}

func GTKMenu_File(menu *gtk.Menu, fpath string, fname string, isdir bool, isapp bool) {
	fpath2 := FolderPathEndSlash(fpath)
	ext_mime := FileMIME(fpath2 + fname)
	app_mime := AppMIME(ext_mime)
	apps_mime := AllAppsMIME(ext_mime)
	GTK_MenuItem(rightmenu, B2S(!isapp, "*", "")+langs.GetStr("cmd_open")+" ["+app_mime+"]", nil)
	if isdir {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_open_new_window"), func() {
			menu.Cancel()
			OpenManager(fpath2 + fname)
		})
		//if(islink){
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_open_symlink"), nil)
	} else {
		GTK_MenuItem(rightmenu, B2S(isapp, "*", "")+langs.GetStr("cmd_app_run"), nil)
		//GTK_MenuSeparator(rightmenu)

		submenu_openwith := GTK_MenuSub(rightmenu, langs.GetStr("cmd_open_with"))
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

		submenu_runas := GTK_MenuSub(rightmenu, langs.GetStr("cmd_app_run_as"))
		GTK_MenuItem(submenu_runas, "SUDO", nil)
		GTK_MenuItem(submenu_runas, "in Terminal", nil)
		GTK_MenuItem(submenu_runas, "SUDO in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "BASH in Terminal", nil)
		//GTK_MenuItem(submenu_runas, "SUDO BASH in Terminal", nil)
	}

	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_file_cut")+" (Ctrl+X)", func() {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fname)
		Prln("cut: " + file1.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{file1}, true)
	})
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_file_copy")+" (Ctrl+C)", func() {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fname)
		Prln("copy: " + file1.GetUrl())
		LinuxClipBoard_CopyFiles([]*LinuxPath{file1}, false)
	})
	GTK_MenuItemIcon(rightmenu, langs.GetStr("cmd_file_delete")+" (Del)", "delete", func() {
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
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_folder_paste_into")+" "+I2S(len(paste_list))+langs.GetStr("cmd_of_objects"), func_paste) // (Ctrl+V)
	}
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_file_rename")+" (F2)", func() {
		Dialog_FileRename(win, fpath2, fname, func() {
			lp := NewLinuxPath(true)
			lp.SetReal(fpath2)
			//listFiles(gGFiles, lp, false)
		})
	})
	if isdir {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_folder_clear"), func() {
			file1 := NewLinuxPath(false) //??
			file1.SetReal(fpath2 + fname)
			Prln("del: " + file1.GetUrl())
			RunFileOperaion([]*LinuxPath{file1}, nil, OPER_CLEAR)
		})
	}
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_files_zip"), nil)
	GTK_MenuSeparator(rightmenu)
	if isapp {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_app_shortcut"), nil)
	}
	if isdir {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_fav_add"), nil)
	}
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_files_info")+" ["+ext_mime+"]", func() {
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
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_selected")+" "+I2S(len(fnames))+" "+langs.GetStr("cmd_of_objects"), nil)
	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_file_cut")+" (Ctrl+X)", func() {
		Prln("cut: " + I2S(len(fnames)) + "files")
		LinuxClipBoard_CopyFiles(list, true)
	})
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_file_copy")+" (Ctrl+C)", func() {
		Prln("copy: " + I2S(len(fnames)) + "files")
		LinuxClipBoard_CopyFiles(list, false)
	})
	GTK_MenuItemIcon(rightmenu, langs.GetStr("cmd_file_delete")+" (Del)", "delete", func() {
		RunFileOperaion(list, nil, OPER_DELETE)
	})
	GTK_MenuSeparator(rightmenu)
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_files_info"), func() {
		Dialog_FileInfo(win, fpath2, fnames)
	})
}

func GTKMenu_FileSearchResult(menu *gtk.Menu, isdir bool, fpath string, fname string) {
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_open_go_folder"), nil)
	if isdir {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_open_new_window"), func() {
			menu.Cancel()
			fpath2 := FolderPathEndSlash(fpath)
			OpenManager(fpath2 + fname)
		})
	} else {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_open_go_folder_of"), func() {
			menu.Cancel()
			fpath2 := FolderPathEndSlash(fpath)
			OpenManager(fpath2)
		})
	}
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_files_info"), func() {
		Dialog_FileInfo(win, FolderPathEndSlash(fpath), []string{fname})
	})
	//GTK_MenuItem(rightmenu, "Info", nil)
}

func GTKMenu_FileSearchResult_Multiple(menu *gtk.Menu, isdir bool, fpath string, fnames []string) {
	GTK_MenuItem(rightmenu, langs.GetStr("cmd_files_info"), func() {
		Dialog_FileInfo(win, FolderPathEndSlash(fpath), fnames)
	})
	//GTK_MenuItem(rightmenu, "Info", nil)
}

func GTKMenu_SMB(menu *gtk.Menu, pc_name string, folder_name string, mounted bool) {
	if StringLength(pc_name) == 0 {
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_smb_login"), nil)
	} else {
		var f_mount func() = nil
		var f_unmount func() = nil
		if !mounted {
			f_mount = func() {
				Dialog_Mount(win, pc_name, folder_name, true)
			}
		} else {
			f_unmount = func() {
				Dialog_Mount(win, pc_name, folder_name, false)
			}
		}
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_smb_mount"), f_mount)
		GTK_MenuItem(rightmenu, langs.GetStr("cmd_smb_unmount"), f_unmount)
	}
	//GTK_MenuSeparator(rightmenu)
}

func GTKMenu_Main(win *gtk.Window) (*gtk.MenuBar, *gtk.Menu) {
	menuBar, _ := gtk.MenuBarNew()
	submenu_file := GTK_MenuSub(menuBar, langs.GetStr("main_cmds"))
	GTK_MenuItem(submenu_file, langs.GetStr("main_new_window"), func() {
		OpenManager(FolderLocation_App())
	})
	GTK_MenuItem(submenu_file, langs.GetStr("main_remote_folder"), nil)
	GTK_MenuItem(submenu_file, langs.GetStr("main_reload_disks"), func() {
		listDiscs(gGDiscs)
	})
	//GTK_MenuItem(submenu_file, "Search", nil)
	/*submenu_view := GTK_MenuSub(menuBar, "View")
	GTK_MenuItem(submenu_view, "Hidden files", nil)
	GTK_MenuItem(submenu_view, "List of files/Icons table", nil)
	GTK_MenuSeparator(submenu_view)
	GTK_MenuItem(submenu_view, "Options", func() {
		Dialog_Options(win)
	})*/

	submenu_edit := GTK_MenuSub(menuBar, langs.GetStr("main_current_folder"))
	//GTKMenu_CurrentFolder(submenu_edit, *path)

	submenu_other := GTK_MenuSub(menuBar, langs.GetStr("main_current_other"))
	GTK_MenuItem(submenu_other, langs.GetStr("main_options"), func() {
		Dialog_Options(win)
	})
	//GTK_MenuItem(submenu_other, "Help", nil)
	GTK_MenuItem(submenu_other, langs.GetStr("main_about"), func() {
		Dialog_About(win, AppVersion(), AppAuthor(), AppMail(), AppRepository(), AppAboutMore(), GetFlag_Russian())
	})
	GTK_MenuItem(submenu_other, langs.GetStr("main_check_update"), func() {
		ok1, v_num1 := GetLastVerison(UrlLastVerison_Manager())
		ok2, v_num2 := GetLastVerison(UrlLastVerison_Mover())
		Prln(B2S_TF(ok1) + v_num1 + "_" + B2S_TF(ok2) + v_num2)
		av1 := AppVersion()
		av2 := AppVersion_Mover(opt.GetFileMover())
		m0 := ""
		m1 := "File Manager installed:\n" + av1
		m2 := "File Mover installed:\n" + av2
		if !ok1 || !ok2 {
			m0 = "Check failed with error(s)\n\n"
		} else {
			if av1 == v_num1 && av2 == v_num2 {
				m0 = "You have last verison of both apps!\n\n"
			} else {
				if av1 != v_num1 {
					m1 = "New version of File Manager!\n" + v_num1 + "\nCurrent verison installed:\n" + av1
				}
				if av2 != v_num2 {
					m2 = "New version of File Mover!\n" + v_num2 + "\nCurrent verison installed:\n" + av2
				}

			}
		}
		Dialog_Message(win, m0+m1+"\n\n"+m2)
	})
	return menuBar, submenu_edit
}
