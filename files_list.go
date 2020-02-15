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
	f               *LinuxFileReport
}

var arr_blocks []*FileIconBlock = []*FileIconBlock{}
var real_files_count = 0
var icon_chanN chan *IconUpdateable
var icon_chan1 chan *IconUpdateable

var async *FileListAsync

func listFiles(g *gtk.Grid, lpath *LinuxPath, scroll_reset bool, save_history bool) {

	//GarbageCollection()

	search, _ := gInpSearch.GetText()
	if save_history {
		hist.SaveNew(lpath.GetVisual(), search)
	}

	upd_title()

	url := lpath.GetUrl()
	is_smb, pc_name, netfolder := SMB_CheckVirtualPath(url)
	if is_smb || (StringLength(pc_name) > 0 && StringLength(netfolder) == 0) {
		listDiscs(gGDiscs)
	}

	if scroll_reset {
		GTK_ScrollReset(sRightScroll)
	}
	if with_destroy {
		for j := 0; j < len(arr_blocks); j++ {
			arr_blocks[j].Destroy()
			arr_blocks[j] = nil
		}
	}
	GTK_Childs(g, true, true)

	spinnerFiles.Start()

	//req :=
	req_id.Add(1)
	fswatcher.Select(lpath.GetReal())

	arr_blocks = []*FileIconBlock{}

	//lpath2 := FolderPathEndSlash(lpath)
	//Prln("=========" + lpath2)

	real_files_count = 0

	rwlock.W_Lock()

	if async != nil {
		async.ForceKill()
	}
	async = NewFileListAsync_DetectType(path, StringTrim(search), 5, 0.2)
	if async == nil {
		Prln("!!! async=nil")
		iconwithlabel := NewFileIconBlock(lpath.GetReal(), "ERROR!", 400, false, false, false, false, "???", ZOOM_SIZE)
		arr_blocks = append(arr_blocks, iconwithlabel)
		g.Attach(iconwithlabel.GetWidgetMain(), 1, 1, 1, 1)
		g.ShowAll()
		return
	}
}

func AddFilesToList(g *gtk.Grid, files []*LinuxFileReport, lpath2 string, req int64) {

	n_now := len(arr_blocks)
	n_max := opt.GetFolderLimit()
	real_files_count += len(files)

	Prln("adding to view " + I2S(len(files)) + " objects")

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

	var arr_render []*IconUpdateable
	_, icon_block_max_w := max_icon_n_w() // icon_block_max_n,icon_block_max_w

	folder_mask := GetIcon_ImageFolder(ZOOM_SIZE)

	for _, f := range files {
		n_now += 1
		if n_now > n_max {
			n_now -= 1
			break
		}
		diff := FolderPathDiff(lpath2, f.FullName.GetReal())
		fname := diff + f.NameOnly
		fname_only := f.NameOnly
		isdir := f.IsDirectory
		isapp := false
		islink := f.IsLink
		isregular := f.IsRegular || islink
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
			inf = inf + FileSizeNiceString(f.SizeBytes) //F2S(float64(f.Size())/float64(BytesInMb), 1) + "Mb"
		}

		ismount := LinuxFolderIsMountPoint(mountlist, lpath2+fname) || SMB_IsMount(path, fname, mountlist)
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
			if f.SizeBytes == 0 {
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
		iconwithlabel.ConnectEventBox("button-press-event", func(_ *gtk.EventBox, event *gdk.Event) {
			disable_focus()
		})
		iconwithlabel.ConnectEventBox("button-release-event", func(_ *gtk.EventBox, event *gdk.Event) {
			//disable_focus()
			mousekey, X, Y, state := GTK_MouseKeyOfEvent(event)
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
						//path.SetReal(path.GetReal() + txtlbl)

						auto_mount := false
						_, pc_name, netfolder := SMB_CheckVirtualPath(path.GetUrl())
						if !ismount && StringLength(pc_name) > 0 && StringLength(netfolder) == 0 {
							auto_mount = true
						}

						path.GoDeeper(txtlbl)
						if opt.GetSymlinkEval() && !path.GetParseProblems() {
							r2, err := FileEvalSymlinks(path.GetReal())
							if err == nil {
								path.SetReal(r2)
							}
						}
						gInpPath.SetText(path.GetVisual())
						gInpSearch.SetText("")

						if auto_mount {
							//Prln("TODO: Add auto-mount here")
							Dialog_Mount(win, pc_name, txtlbl, true)
						}

						listFiles(gGFiles, path, true, true)
					} else {
						OpenFileByApp(path.GetReal()+txtlbl, "")
					}
				} else {
					clicktime = TimeNow()
					if X > 20 || Y > 20 {
						Prln(">>click at file block")
						if !GTK_KeyboardCtrlState(state) {
							FilesSelector_ResetChecks()
							iconwithlabel.SetSelected(true)
						} else {
							iconwithlabel.SetSelected(!iconwithlabel.GetSelected())
						}
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
				if diff == "" {
					url := path.GetUrl()
					is_smb, pc_name, netfolder := SMB_CheckVirtualPath(url)
					if is_smb || (StringLength(pc_name) > 0 && StringLength(netfolder) == 0) {
						GTKMenu_SMB(rightmenu, pc_name, fname, ismount)
					} else {
						if len(sel_list) <= 1 {
							GTKMenu_File(rightmenu, lpath2, fname, isdir, isapp)
						} else {
							GTKMenu_Files(rightmenu, lpath2, sel_list, isdir, isapp)
						}
					}
				} else {
					if len(sel_list) <= 1 {
						GTKMenu_FileSearchResult(rightmenu, isdir, lpath2+diff, fname_only)
					} else {
						GTKMenu_FileSearchResult_Multiple(rightmenu, isdir, lpath2, sel_list)
					}
				}

				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event) // (evBox, gdk.GDK_GRAVITY_STATIC, gdk.GDK_GRAVITY_STATIC,
			}
		})

		arr_blocks = append(arr_blocks, iconwithlabel)

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
						arr_render[j].pixbuf_cache = CachePreview_ReadPixbuf(arr_render[j].f.FileReport(), ZOOM_SIZE, folder_mask)
						if arr_render[j].pixbuf_cache != nil {
							qu.Append(arr_render[j])
						}
					}
				} else {
					if path.GetReal() != opt.GetHashFolder() {
						arr_render[j].pixbuf_cache = CachePreview_ReadPixbuf(arr_render[j].f.FileReport(), ZOOM_SIZE, nil)
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
		rwlock.R_Lock()
		for j := 0; j < len(arr_render); j++ {
			if req == req_id.Get() {
				arr_render[j].req = req
				if single_thread_protocol {
					icon_chan1 <- arr_render[j]
				} else {
					icon_chanN <- arr_render[j]
				}
			}
			RuntimeGosched()
		}
		rwlock.R_Unlock()
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

func FolderPathDiff(origpath string, fullname string) string {
	orarr := StringSplit(FilePathEndSlashRemove(origpath), "/")
	len1 := len(orarr)
	farr := StringSplit(FilePathEndSlashRemove(fullname), "/")
	len2 := len(farr) - 1
	if len2 > len1 {
		//Prln("??" + origpath + " >> " + fullname)
		return FolderPathEndSlash(StringJoin(farr[len1:len2], "/"))
	}
	return ""
}

type PathHistory struct {
	pathes   []string
	searches []string
	ind      int
}

func PathHistoryNew() *PathHistory {
	h := PathHistory{pathes: []string{}, searches: []string{}, ind: 0}
	return &h
}

func (h *PathHistory) SaveNew(path_visual string, search string) {
	if h.ind > 0 {
		h.pathes = h.pathes[0:h.ind]
		h.searches = h.searches[0:h.ind]
	} else {
		h.pathes = []string{}
		h.searches = []string{}
	}

	L := len(h.pathes)
	if L > 0 && h.pathes[L-1] == path_visual && h.searches[L-1] == search {
		return
	}

	h.pathes = append(h.pathes, path_visual)
	h.searches = append(h.searches, search)

	h.ind++
}

func (h *PathHistory) Back() (bool, string, string) {
	if h.ind > 1 {
		h.ind--
		return true, h.pathes[h.ind-1], h.searches[h.ind-1]
	} else {
		return false, "", ""
	}
}

func (h *PathHistory) Forward() (bool, string, string) {
	if h.ind < len(h.pathes) {
		h.ind++
		return true, h.pathes[h.ind-1], h.searches[h.ind-1]
	} else {
		return false, "", ""
	}
}

func (h *PathHistory) CanBackForward() (bool, bool) {
	return h.ind > 1, h.ind < len(h.pathes)
}

func (h *PathHistory) GetList() (int, []string, []string) {
	return h.ind, h.pathes, h.searches
}

func (h *PathHistory) GoAt(id int) (bool, string, string) {
	//TODO
	return false, "", ""
}

func (h *PathHistory) GetCurrent() (string, string) {
	i := h.ind - 1
	if len(h.pathes) > 0 && i >= 0 && i < len(h.pathes) {
		return h.pathes[i], h.searches[i]
	} else {
		Prln("History GetCurrent(): Error")
		return "", ""
	}

}
