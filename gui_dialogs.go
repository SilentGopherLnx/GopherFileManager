package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"github.com/gotk3/gotk3/gtk"
)

func Dialog_FileRename(w *gtk.Window, fpath string, fname_old string, upd func()) {
	dial, box := GTK_Dialog(w, "Rename: "+fname_old)
	dial.SetDefaultSize(350, 90)

	inpname, _ := gtk.EntryNew()
	inpname.SetText(fname_old)
	inpname.SetHExpand(true)
	lbl_err, _ := gtk.LabelNew("")
	lbl_err.SetHExpand(true)
	lbl_err.SetVExpand(true)
	btnok, _ := gtk.ButtonNewWithLabel("Ok")
	btnok.SetHExpand(true)
	btnok.Connect("button-press-event", func() {
		safe_name, _ := inpname.GetText()
		// Windows (FAT32, NTFS): Any Unicode except NUL, \, /, :, *, ", <, >, |
		// Mac(HFS, HFS+): Any valid Unicode except : or /
		// Linux(ext[2-4]): Any byte except NUL or /
		if safe_name != fname_old {
			fpath2 := FolderPathEndSlash(fpath)
			ok := FileRename(fpath2+fname_old, fpath2+safe_name)
			if ok {
				dial.Close()
				upd()
			} else {
				lbl_err.SetText("error")
			}
		} else {
			dial.Close()
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
	dial.Run()
	dial.Close()
}

func Dialog_FileInfo(w *gtk.Window, fpath string, fname string) {
	fullname := FolderPathEndSlash(fpath) + fname
	winw, winh := 300, 300
	win2, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err == nil {
		win2.SetTitle("Info: " + fname)
		win2.SetDefaultSize(winw, winh)
		win2.SetPosition(gtk.WIN_POS_CENTER)
		//win2.SetTransientFor(w)
		win2.SetIconFromFile(FolderLocation_App() + "gui/icon.png")
		//win2.SetModal(true)
		//win2.SetKeepAbove(true)

		src_size := NewAtomicInt64(0)
		src_files := NewAtomicInt64(0)
		src_folders := NewAtomicInt64(0)
		src_failed := NewAtomicInt64(0)
		src_irregular := NewAtomicInt64(0)
		src_mount := NewAtomicInt64(0)
		src_symlinks := NewAtomicInt64(0)

		box_size, lbl_size := GTK_LabelPair("Size: ", "calculating...")
		box_filse, lbl_files := GTK_LabelPair("Objects: ", "calculating...")
		part_name, _ := LinuxFilePartition(mountlist, fullname)
		box_disk, _ := GTK_LabelPair("Disk: ", part_name)
		box_mime, _ := GTK_LabelPair("Mime type: ", FileMIME(fullname))
		box_perm, _ := GTK_LabelPair("Permissions: ", FilePermissionsString(fullname))
		go func() {
			path_info := FolderPathEndSlash(fpath) + fname
			file_or_dir, ok := FileInfo(path_info)
			if ok {
				FoldersRecursively_Size(mountlist, file_or_dir, path_info, src_size, src_files, src_folders, src_failed, src_irregular, src_mount, src_symlinks)
			}
		}()
		grid, _ := gtk.GridNew()
		grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
		grid.Add(box_size)
		grid.Add(box_filse)
		grid.Add(box_disk)
		grid.Add(box_mime)
		grid.Add(box_perm)
		win2.Add(grid)

		upd_info := func() {
			sel_size := src_size.Get()
			lbl_size.SetText(FileSizeNiceString(sel_size) + " (" + I2Ss(sel_size) + " bytes)")
			sel_files := src_files.Get()
			sel_folders := src_folders.Get()
			sel_failed := src_failed.Get()
			sel_irregular := src_irregular.Get()
			sel_mount := src_mount.Get()
			sel_symlinks := src_symlinks.Get()
			lbl_files.SetText(I2S64(sel_files+sel_folders) + " (" + I2S64(sel_files) + " files & " + I2S64(sel_folders) + " folders) " +
				"\n" + I2S64(sel_failed) + " folders content is blocked for reading" +
				"\n+ " + I2S64(sel_irregular) + " irregular files, " + I2S64(sel_mount) + " mount points, " + I2S64(sel_symlinks) + " symlinks")
		}
		win2.Connect("destroy", func() {
			main_iterations_funcs.Remove(&upd_info)
		})

		win2.ShowAll()
		win2.SetSizeRequest(winw, winh)
		main_iterations_funcs.Add(&upd_info)
	} else {
		Prln(err.Error())
	}
}
