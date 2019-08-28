package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "./pkg_fileicon"
	. "./pkg_filetools"

	//	"os/exec"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/glib"
)

type IconUpdateable struct {
	widget         *FileIconBlock
	pixbuf_cache   *gdk.Pixbuf // loaded old cache of preview
	pixbuf_preview *gdk.Pixbuf // loaded preview
	fullname       string
	fname          string
	tfile          string
	basic_mode     bool //check is executable
	folder         bool
	//oldbuf         bool //have loaded old preview
	success_preview bool
	req             int64
	f               FileReport
}

func listFiles(g *gtk.Grid, lpath string, scroll_reset bool) {

	req := req_id.Add(1)

	if scroll_reset {
		GTK_ScrollReset(sRightScroll)
	}

	fswatcher.Select(lpath)

	new_ind := path_updated.Add(1)

	if with_destroy {
		for j := 0; j < len(arr_blocks); j++ {
			arr_blocks[j].Destroy()
			arr_blocks[j] = nil
		}
	}

	arr_blocks = []*FileIconBlock{}

	//GarbageCollection()

	GTK_Childs(g, true, true)

	GarbageCollection()

	lpath2 := FolderPathEndSlash(lpath)
	Prln("=========" + lpath2)

	single_thread_protocol := false
	with_extra_info := true
	if StringFind(lpath2, "/run/user/") == 1 {
		single_thread_protocol = true
		if StringFind(lpath2, "/gvfs/smb-share:") > 1 {
			single_thread_protocol = false
		} else {
			Prln("single_thread_protocol TRUE")
		}
		with_extra_info = false
	}
	if path.GetReal() == opt.GetHashFolder() {
		with_extra_info = false
	}

	files, err := Folder_ListFiles(lpath2, false) //// !!!!!!!!!!!!!!!!!!!!!!! true!
	if err != nil {
		Prln(err.Error())
		iconwithlabel := NewFileIconBlock(lpath2, "ERROR!", 400, false, false, false, false, err.Error(), ZOOM_SIZE)
		arr_blocks = append(arr_blocks, iconwithlabel)
		g.Attach(iconwithlabel.GetWidgetMain(), 1, 1, 1, 1)
		g.ShowAll()
		return
	}
	j := 0

	var arr_render []*IconUpdateable
	_, icon_block_max_w := max_icon_n_w() // icon_block_max_n,icon_block_max_w

	folder_mask := GetIcon_ImageFolder(ZOOM_SIZE)

	for _, f := range files {
		fname := f.Name()
		isdir := f.IsDir()
		isapp := false
		islink := FileIsLink(f)
		isregular := f.Mode().IsRegular() || islink
		//oldbuf := false

		if islink {
			isdir = FileLinkIsDir(lpath2 + fname)
		}

		filepathfinal := lpath2 + fname
		if isdir {
			filepathfinal = FolderPathEndSlash(filepathfinal)
		}

		//Prln("[" + B2S_YN(isdir) + "]:{" + fname + "}" + B2S_YN(islink) + "/" + f.Mode().String())

		inf := "" //f.Mode().String() + "\n" // + "|" + f.Mode().Perm().String()

		not_read := false
		if isdir {
			if !single_thread_protocol && with_extra_info {
				fl, err := Folder_ListFiles(filepathfinal, false)
				if err == nil {
					if !single_thread_protocol {
						inf = inf + I2S(len(fl)) + " files"
					}
				} else {
					not_read = true
				}
			}
		} else {
			inf = inf + FileSizeNiceString(f.Size()) //F2S(float64(f.Size())/float64(BytesInMb), 1) + "Mb"
		}

		ismount := LinuxFolderIsMountPoint(mountlist, lpath2+fname)
		iconwithlabel := NewFileIconBlock(lpath2, fname, icon_block_max_w, isdir, islink, not_read, ismount, inf, ZOOM_SIZE)

		if isdir {
			if filepathfinal == opt.GetHashFolder() {
				iconwithlabel.SetIconPixPuf(GetIcon_PixBif_OF(ZOOM_SIZE, PREFIX_DRAWONME+FILE_TYPE_FOLDER_HASH))
			} else {
				iconwithlabel.SetIconPixPuf(GetIcon_PixBif(ZOOM_SIZE, "", true))
			}
			dest := NewLinuxPath(true)
			dest.SetReal(lpath2 + fname)
			GTK_CopyPasteDnd_SetFolderDest(iconwithlabel.GetWidgetMain(), dest)
		} else {
			tfile := FileExtension(fname)
			var pixbuf_icon *gdk.Pixbuf = nil
			mime := ""
			if with_extra_info {
				mime = FileMIME(filepathfinal)
			}
			if mime == APP_EXEC_TYPE {
				isapp = true
				pixbuf_icon = GetIcon_PixBif_OF(ZOOM_SIZE, PREFIX_EXTRA+FILE_TYPE_BIN)
			} else {
				pixbuf_icon = GetIcon_PixBif(ZOOM_SIZE, tfile, false)
			}
			if f.Size() == 0 {
				if isregular {
					pixbuf_icon = GetIcon_PixBif_OF(ZOOM_SIZE, PREFIX_DRAWONME+FILE_TYPE_ZERO)
				} else {
					pixbuf_icon = GetIcon_PixBif_OF(ZOOM_SIZE, PREFIX_DRAWONME+FILE_TYPE_NOTFILE)
				}
				//Prln("zero")
			}
			iconwithlabel.SetIconPixPuf(pixbuf_icon)
		}

		srcfile := NewLinuxPath(isdir)
		srcfile.SetReal(lpath2 + fname)
		getter := func() []*LinuxPath {
			list := []*LinuxPath{}
			fnames := FilesSelector_GetList()
			if len(fnames) > 1 {
				for j := 0; j < len(fnames); j++ {
					file1 := NewLinuxPath(false) //??
					file1.SetReal(lpath2 + fnames[j])
					list = append(list, file1)
				}
			} else {
				list = []*LinuxPath{srcfile}
			}
			return list
		}
		GTK_CopyPasteDnd_SetIconSource(iconwithlabel.GetWidgetMain(), iconwithlabel.GetIcon(), getter)

		clicktime := TimeAddMS(TimeNow(), -2000)
		iconwithlabel.ConnectEventBox("button-release-event", func(_ *gtk.EventBox, event *gdk.Event) {
			mousekey, X, Y := GTK_MouseKeyOfEvent(event)
			switch mousekey {
			case 1:
				dt := AbcF(TimeSeconds(clicktime))
				//Prln(I2S(dt) + " / " + TimeStr(clicktime))
				if dt < 0.5 {
					txtlbl := iconwithlabel.GetFileName()
					Prln("click: [" + txtlbl + "]")
					clicktime = TimeAddMS(clicktime, -2000)
					if isdir {
						//path, _ = gInpPath.GetText()
						path.SetReal(path.GetReal() + txtlbl)
						if opt.GetSymlinkEval() {
							r2, err := FileEvalSymlinks(path.GetReal())
							if err == nil {
								path.SetReal(r2)
							}
						}
						gInpPath.SetText(path.GetVisual())
						listFiles(gGFiles, path.GetReal(), true)
					} else {
						OpenFileByApp(path.GetReal()+txtlbl, "")
					}
				} else {
					clicktime = TimeNow()
					if X > 20 || Y > 20 {
						Prln(">>click at file block")
						FilesSelector_ResetChecks()
						iconwithlabel.SetSelected(true)
					}
				}
				//gGFiles.QueueDraw()
			case 3:
				Prln("right")
				if rightmenu != nil && rightmenu.IsVisible() {
					Prln("hiding menu")
					rightmenu.Destroy()
				}
				rightmenu, _ = gtk.MenuNew()

				sel_list := FilesSelector_GetList()
				fname2 := fname
				if isdir {
					fname2 = FolderPathEndSlash(fname)
				}
				if StringInArray(fname2, sel_list) == -1 {
					FilesSelector_ResetChecks()
					sel_list = []string{}
				}
				if len(sel_list) <= 1 {
					GTKMenu_File(rightmenu, lpath2, fname, isdir, isapp)
				} else {
					GTKMenu_Files(rightmenu, lpath2, sel_list, isdir, isapp)
				}

				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event) // (evBox, gdk.GDK_GRAVITY_STATIC, gdk.GDK_GRAVITY_STATIC,
			}
		})

		arr_blocks = append(arr_blocks, iconwithlabel)
		j++

		if isdir {
			fullname := FolderPathEndSlash(path.GetReal() + fname)
			if with_folders_preview && !single_thread_protocol && fullname != opt.GetHashFolder() {
				iconwithlabel.SetLoading(true, false)
				arr_render = append(arr_render, &IconUpdateable{widget: iconwithlabel, fullname: fullname, fname: fname, tfile: "", basic_mode: single_thread_protocol, folder: true, f: f})
			}
		} else {
			/*if !single_thread_protocol {
				update_icon(lpath2, fname, icon)
			}*/
			if with_files_preview {
				tfile := FileExtension(fname)
				if len(tfile) > 0 {
					fullname := path.GetReal() + fname
					if FileIsPreviewAbble(tfile) && path.GetReal() != opt.GetHashFolder() {
						iconwithlabel.SetLoading(true, false)
						arr_render = append(arr_render, &IconUpdateable{widget: iconwithlabel, fullname: fullname, fname: fname, tfile: tfile, basic_mode: single_thread_protocol, folder: false, f: f})
					}
				}
			}
		}
	}

	resort_and_show()
	g.ShowAll()
	//win.QueueDraw()

	//Prln("folder loaded. starting chans sending...")

	go func() {
		SleepMS(5)

		if with_cache_preview {
			SortArray(arr_render, func(i, j int) bool {
				if arr_render[i].folder != arr_render[j].folder {
					return !CompareBoolLess(arr_render[i].folder, arr_render[j].folder) // first folders
				}
				return false
			})

			for j := 0; j < len(arr_render); j++ {
				if arr_render[j].folder {
					if FolderPathEndSlash(arr_render[j].fullname) != opt.GetHashFolder() {
						arr_render[j].pixbuf_cache = CachePreview_ReadPixbuf(arr_render[j].f, ZOOM_SIZE, folder_mask)
						if arr_render[j].pixbuf_cache != nil {
							qu.Append(arr_render[j])
						}
					}
				} else {
					if path.GetReal() != opt.GetHashFolder() {
						arr_render[j].pixbuf_cache = CachePreview_ReadPixbuf(arr_render[j].f, ZOOM_SIZE, nil)
						if arr_render[j].pixbuf_cache != nil {
							qu.Append(arr_render[j])
						}
					}
				}
			}
		}

		SortArray(arr_render, func(i, j int) bool {
			pixbuf_cache1 := arr_render[i].pixbuf_cache != nil // loaded
			pixbuf_cache2 := arr_render[j].pixbuf_cache != nil
			if pixbuf_cache1 != pixbuf_cache2 {
				return CompareBoolLess(pixbuf_cache1, pixbuf_cache2) // first not loaded
			}
			if arr_render[i].folder != arr_render[j].folder {
				if pixbuf_cache1 { // loaded
					return !CompareBoolLess(arr_render[i].folder, arr_render[j].folder) // first folders (other threads can load files quicker if num folders < num threads)
				} else {
					return CompareBoolLess(arr_render[i].folder, arr_render[j].folder) // first files
				}
			}
			return false
		})
		for j := 0; j < len(arr_render); j++ {
			if new_ind == path_updated.Get() {
				arr_render[j].req = req
				if single_thread_protocol {
					icon_chan1 <- arr_render[j]
				} else {
					icon_chanN <- arr_render[j]
				}
			}
			RuntimeGosched()
		}

		//Prln("GO FINISH")
	}()

}

func resort_and_show() {
	SortArray(arr_blocks, func(i, j int) bool {
		if arr_blocks[i].IsDir() != arr_blocks[j].IsDir() {
			return !CompareBoolLess(arr_blocks[i].IsDir(), arr_blocks[j].IsDir())
		}
		name1 := arr_blocks[i].GetFileName()
		name2 := arr_blocks[j].GetFileName()
		if sort_mode == 1 {
			type1 := FileExtension(name1)
			type2 := FileExtension(name2)
			if type1 != type2 {
				return XOR(type1 < type2, sort_reverse)
			}
		}
		if name1 != name2 {
			return XOR(FileSortName(StringDown(name1)) < FileSortName(StringDown(name2)), sort_reverse)
		}
		return false
	})

	// ============
	icon_block_max_n, icon_block_max_w := max_icon_n_w()

	icon_block_max_n_old = icon_block_max_n
	icon_block_max_w_old = icon_block_max_w
	icon_block_max_w += 0

	GTK_Childs(gGFiles, true, false)

	for j := 0; j < len(arr_blocks); j++ {
		x := j % icon_block_max_n
		y := j / icon_block_max_n
		gEv := arr_blocks[j].GetWidgetMain()
		gGFiles.Attach(gEv, x, y, 1, 1)
		arr_blocks[j].SetWidth(icon_block_max_w)
	}
}

func resize_event_no_repeats() {
	icon_block_max_n, icon_block_max_w := max_icon_n_w()
	if icon_block_max_n_old != icon_block_max_n || icon_block_max_w_old != icon_block_max_w {
		Prln("resized")
		resort_and_show()
	}
}

func max_icon_n_w() (int, int) {
	//ww, _ := sScroll.GetPreferredWidth()
	//sScroll.CheckResize()
	//ww, _ := win.GetPreferredWidth()
	ww, _ := win.GetSize()
	real_w := MAXI(16, ww-LEFT_PANEL_SIZE) - 6 - BORDER_SIZE*5/2
	icon_block_max_w := MAXI(16, ZOOM_SIZE+BORDER_SIZE*4)
	icon_block_max_n := MAXI(1, MAXI(16, real_w)/icon_block_max_w)
	icon_block_max_w = real_w/icon_block_max_n - BORDER_SIZE*3
	//Prln("size" + I2S(ww))
	return icon_block_max_n, icon_block_max_w
}
