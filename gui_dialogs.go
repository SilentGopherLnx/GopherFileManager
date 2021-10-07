package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_fileicon"
	. "github.com/SilentGopherLnx/GopherFileManager/pkg_filetools"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func Dialog_FileRename(w *gtk.Window, fpath string, fname_old string, upd func()) {
	if opt.GetRenameByThisApp() {
		Dialog_FileRename_ThisApp(w, fpath, fname_old, upd)
	} else {
		Dialog_FileRename_ExternalApp(w, fpath, fname_old, upd)
	}
}

func Dialog_FileRename_ExternalApp(w *gtk.Window, fpath string, fname_old string, upd func()) {
	file1 := NewLinuxPath(false) //??
	file1.SetReal(FolderPathEndSlash(fpath) + fname_old)
	Prln("rename: " + file1.GetUrl())
	RunFileOperaion([]*LinuxPath{file1}, nil, OPER_RENAME)
}

func Dialog_FileRename_ThisApp(w *gtk.Window, fpath string, fname_old string, upd func()) {
	fname_prev := FilePathEndSlashRemove(fname_old)
	dial, box := GTK_Dialog(w, langs.GetStr("dialog_rename_title")+": "+fname_prev)
	dial.SetDefaultSize(350, 90)

	inpname, _ := gtk.EntryNew()
	inpname.SetText(fname_prev)
	inpname.SetHExpand(true)
	lbl_err, _ := gtk.LabelNew("")
	lbl_err.SetHExpand(true)
	lbl_err.SetVExpand(true)
	GTK_LabelWrapMode(lbl_err, 1)
	btnok, _ := gtk.ButtonNewWithLabel("Ok")
	btnok.SetHExpand(true)

	//var close_dial *ABool = NewAtomicBool(false)

	ok_func := func() {
		safe_name, _ := inpname.GetText()
		// Windows (FAT32, NTFS): Any Unicode except NUL, \, /, :, *, ", <, >, |
		// Mac(HFS, HFS+): Any valid Unicode except : or /
		// Linux(ext[2-4]): Any byte except NUL or /
		if safe_name != fname_prev {
			fpath2 := FolderPathEndSlash(fpath)
			ok, errtxt := FileRename(fpath2+fname_prev, fpath2+safe_name)
			if ok {
				//close_dial.Set(true)
				//dial.Close()
				dial.Destroy()
				if upd != nil {
					upd()
				}
			} else {
				lbl_err.SetText(langs.GetStr("msg_error") + ": " + errtxt)
				//close_dial.Set(false)
			}
		} else {
			//close_dial.Set(true)
			//dial.Close()
			dial.Destroy()
		}
	}

	btnok.Connect("button-press-event", ok_func)
	dial.Connect("key-press-event", func(_ *gtk.Dialog, ev *gdk.Event) {
		uint_key, _, _ := GTK_KeyboardKeyOfEvent(ev)
		if uint_key == GTK_KEY_Enter() {
			ok_func()
		}
	})

	box.SetOrientation(gtk.ORIENTATION_VERTICAL)
	box.Add(inpname)
	box.Add(lbl_err)
	box.Add(btnok)

	// box.SetSpacing(0)
	// box.SetBorderWidth(0)
	// box.SetMarginBottom(0)

	// dial.SetBorderWidth(0)
	// dial.SetMarginBottom(0)

	dial.ShowAll()
	ind := StringFindEnd(fname_prev, ".")
	if ind == 0 || ind == 1 {
		ind = StringLength(fname_prev)
	} else {
		ind--
	}
	inpname.SelectRegion(0, ind)
	//if !close_dial.Get() {
	//dial.Run()
	//Prln("!+!+")
	//dial.Destroy()
	//}
	//dial.Close()

}

func Dialog_FileInfo(w *gtk.Window, fpath string, fnames []string) {
	winw, winh := 300, 300
	win2, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	//over := NewAtomicBool(false, [2]string{"1", "0"})
	if err == nil {
		kill := NewAtomicBool(false)
		if len(fnames) == 1 {
			win2.SetTitle(langs.GetStr("dialog_fileinfo_title") + ": \"" + fnames[0] + "\"")
		} else {
			win2.SetTitle(langs.GetStr("dialog_fileinfo_title") + ": " + I2S(len(fnames)) + " " + langs.GetStr("dialog_fileinfo_title_selected"))
		}
		win2.SetDefaultSize(winw, winh)
		win2.SetPosition(gtk.WIN_POS_CENTER)
		//win2.SetTransientFor(w)
		win2.SetIconFromFile(FolderLocation_App() + GUI_PATH + "icon.png")
		//win2.SetModal(true)
		//win2.SetKeepAbove(true)

		src_size := NewAtomicInt64(0)
		src_files := NewAtomicInt64(0)
		src_folders := NewAtomicInt64(0)
		src_failed := NewAtomicInt64(0)
		src_irregular := NewAtomicInt64(0)
		src_mount := NewAtomicInt64(0)
		src_symlinks := NewAtomicInt64(0)

		spinner, _ := gtk.SpinnerNew()
		spinner.Start()

		mime := "?"
		perm := "?"
		if len(fnames) == 1 {
			fullname := FolderPathEndSlash(fpath) + fnames[0]
			mime = FileMIME(fullname)
			perm = FilePermissionsString(fullname)
		}

		box_path, _ := GTK_LabelPair(langs.GetStr("dialog_fileinfo_path")+": ", fpath)

		names_str := ""
		for j := 0; j < len(fnames); j++ {
			names_str += fnames[j] + "\n"
		}

		text_selected_files := langs.GetStr("dialog_fileinfo_selected") + ":"
		lbl_src_title, _ := gtk.LabelNew(text_selected_files)
		lbl_src_title.SetMarkup("<b>" + text_selected_files + "</b>")

		lbl_src, _ := gtk.LabelNew(names_str)
		lbl_src.SetHExpand(true)
		lbl_src.SetVAlign(gtk.ALIGN_START)
		lbl_src.SetHAlign(gtk.ALIGN_START)
		//lbl_src.SetJustify(gtk.JUSTIFY_LEFT)
		GTK_LabelWrapMode(lbl_src, MAXI(1, len(fnames)))

		scroll_scr, _ := gtk.ScrolledWindowNew(nil, nil)
		//scroll_scr.SetVExpand(true)
		scroll_scr.SetHExpand(true)
		scroll_scr.Add(lbl_src)
		scroll_scr.SetSizeRequest(0, 50)

		frame, _ := gtk.FrameNew(text_selected_files)
		frame.SetLabelWidget(lbl_src_title)
		frame.Add(scroll_scr)

		box_size, lbl_size := GTK_LabelPair(langs.GetStr("dialog_fileinfo_size")+": ", langs.GetStr("dialog_fileinfo_calculating"))
		box_filse, lbl_files := GTK_LabelPair(langs.GetStr("dialog_fileinfo_objects")+": ", langs.GetStr("dialog_fileinfo_calculating"))
		part_name, _ := LinuxFilePartition(mountlist, fpath)
		box_disk, _ := GTK_LabelPair(langs.GetStr("dialog_fileinfo_disk")+": ", part_name)
		box_mime, _ := GTK_LabelPair(langs.GetStr("dialog_fileinfo_mime")+": ", mime)
		box_perm, _ := GTK_LabelPair(langs.GetStr("dialog_fileinfo_permissions")+": ", perm)
		go func() {
			for j := 0; j < len(fnames); j++ {
				path_info := FolderPathEndSlash(fpath) + fnames[j]
				file_or_dir, err := FileInfo(path_info, false)
				if err == nil {
					FoldersRecursively_Size(mountlist, file_or_dir, path_info, src_size, src_files, src_folders, src_failed, src_irregular, src_mount, src_symlinks, kill)
				}
			}
			//over.Set(true)
			spinner.Stop()
		}()
		grid, _ := gtk.GridNew()
		grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
		grid.Add(box_path)
		grid.Add(frame)
		grid.Add(spinner)
		grid.Add(box_size)
		grid.Add(box_filse)
		grid.Add(box_disk)
		grid.Add(box_mime)
		grid.Add(box_perm)
		win2.Add(grid)

		upd_info := func() {
			sel_size := src_size.Get()
			lbl_size.SetText(FileSizeNiceString(sel_size) + " (" + I2Ss(sel_size) + " " + langs.GetStr("dialog_fileinfo_bytes") + ")")
			sel_files := src_files.Get()
			sel_folders := src_folders.Get()
			sel_failed := src_failed.Get()
			sel_irregular := src_irregular.Get()
			sel_mount := src_mount.Get()
			sel_symlinks := src_symlinks.Get()
			lbl_files.SetText(I2S64(sel_files+sel_folders) +
				" (" + I2S64(sel_files) + " " + langs.GetStr("dialog_fileinfo_files") + " & " + I2S64(sel_folders) + " " + langs.GetStr("dialog_fileinfo_folders") + ") " +
				"\n" + I2S64(sel_failed) + " folders content is blocked for reading" +
				"\n+ " + I2S64(sel_irregular) + " irregular files, " + I2S64(sel_mount) + " mount points, " + I2S64(sel_symlinks) + " symlinks")
		}
		win2.Connect("destroy", func() {
			main_iterations_funcs.Remove(&upd_info)
			kill.Set(true)
		})

		win2.ShowAll()
		win2.SetSizeRequest(winw, winh)
		main_iterations_funcs.Add(&upd_info)
	} else {
		Prln(err.Error())
	}
}

func Dialog_FolderError(w *gtk.Window, err error, path_visual string) {
	Prln("ERROR DIALOG")
	dial, box := GTK_Dialog(w, langs.GetStr("dialog_error_title"))

	lbl_err, _ := gtk.LabelNew(StringFill(err.Error(), 20))

	box.SetOrientation(gtk.ORIENTATION_VERTICAL)
	box.Add(lbl_err)

	dial.SetResizable(false)
	dial.ShowAll()
	dial.Run()
	dial.Close()
}

func Dialog_Message(w *gtk.Window, text string) {
	dial, box := GTK_Dialog(w, langs.GetStr("dialog_message_title"))

	lbl_err, _ := gtk.LabelNew(StringFill(text, 20))

	box.SetOrientation(gtk.ORIENTATION_VERTICAL)
	box.Add(lbl_err)

	dial.SetResizable(false)
	dial.ShowAll()
	dial.Run()
	dial.Close()
}

func Dialog_Mount(w *gtk.Window, pc_name string, folder_name string, mount bool) {
	if !mount {
		go func() {
			err := SMB_UnMount(pc_name, folder_name)
			if err != nil {
				Prln(err.Error())
			}
			fswatcher.EmulateUpdate()
		}()
		return
	}

	go func() {
		_, err := SMB_MountLoginAsk(pc_name, folder_name)
		if err != nil {
			Prln(err.Error())
		}
		fswatcher.EmulateUpdate()
	}()

	dial, box := GTK_Dialog(w, langs.GetStr("dialog_mounting_title"))

	lbl, _ := gtk.LabelNew(StringFill("Wait some time and close this window.\nUse default file manager if you have password and it's not saved in system", 20))

	box.SetOrientation(gtk.ORIENTATION_VERTICAL)
	box.Add(lbl)

	dial.SetResizable(false)
	dial.ShowAll()
	dial.Run()
	dial.Close()
}
