package main

//sudo apt-get install libgtk-3-dev
//sudo apt-get install libcairo2-dev
//sudo apt-get install libglib2.0-dev
//https://github.com/gotk3/gotk3
//https://github.com/golang/image

//sudo apt install ffmpeg

//xclip

//ghex

//godoc -http=":6060"

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// #cgo CFLAGS: -Wno-deprecated
//import "C"

var win *gtk.Window

var gGCenter *gtk.Paned
var sLeftScroll, sRightScroll *gtk.ScrolledWindow
var gGFiles *gtk.Grid
var gGDiscs *gtk.Box
var gInpPath *gtk.Entry
var gBtnUp *gtk.Button
var gBtnRefresh *gtk.Button
var mem, space *gtk.Label

var path *LinuxPath = NewLinuxPath(true)
var path_updated = NewAtomicInt(0)

var icon_block_max_n_old, icon_block_max_w_old int

var ZOOM_SIZE = 64 * 2

var LEFT_PANEL_SIZE = 200 //200

var BORDER_SIZE = 8

var arr_blocks []*GtkFileIconBlock = []*GtkFileIconBlock{}
var icon_chanN chan *IconUpdateable
var icon_chan1 chan *IconUpdateable

var qu *SyncQueue

var with_folders_preview bool = false
var with_files_preview bool = false

var with_destroy bool = true

//var with_extra bool = false

var usage = ""

var rightmenu *gtk.Menu = nil

var main_iterations_funcs *FuncArr = NewFuncArr()

var mountlist [][2]string

var fswatcher *FSWatcher

var num_works *AInt = NewAtomicInt(0)

func init() {

	AboutVersion(AppVersion())

	//Prln(StringFill("123456", 5))

	//TestLinuxPath()

	// inf, _ := FileInfo("/home/nike/.steam/root/")
	// Prln(B2S_YN(FileIsLink(inf)))

	InitOptions()

	ZOOM_SIZE = opt.GetZoom()

	args := AppRunArgs()
	if len(args) >= 2 {
		path.SetReal(args[1])
	} else {
		path.SetReal(FolderLocation_WorkDir())
	}

	icon_chanN = make(chan *IconUpdateable)
	icon_chan1 = make(chan *IconUpdateable)

	qu = NewSyncQueue()

	//path = "/mnt/dm-1/"
	with_folders_preview = true
	with_files_preview = true
	//with_destroy = false

}

func main() {

	fswatcher = NewFSWatcher()
	defer fswatcher.Close()

	gtk.Init(nil)

	var err error
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		Prln("Unable to create window:") // + err)
	}
	uid, _, _ := GetPC_UserUidLoginName()
	sudo := Select_String(LinuxRoot_Check() == 1, "[root"+uid+"] ", "")
	win.SetTitle(sudo + "GopherFileManager")
	win.SetDefaultSize(1200, 800)
	win.SetPosition(gtk.WIN_POS_CENTER)

	win.Connect("destroy", func() {
		AppExit(0)
	})

	//ev := 0 // https://developer.gnome.org/gtk3/stable/GtkWidget.html
	win.Connect("size-allocate", func() {
		resize_files_icons()
	})

	win.SetIconFromFile(appdir + "gui/icon.png")

	// ================

	gBtnUp, _ = gtk.ButtonNewWithLabel("Up")
	//gBtnUp.SetProperty("background-color", "red")
	//img1 := GTK_Image_From_File(appdir+"gui/button_up.png", "png")
	img1 := GTK_Image_From_Name("go-up", gtk.ICON_SIZE_BUTTON)
	gBtnUp.SetImage(img1)
	gBtnUp.SetProperty("always-show-image", true)
	gBtnUp.Connect("clicked", func() {
		path.GoUp()
		gInpPath.SetText(path.GetVisual())
		listFiles(gGFiles, path.GetReal())
	})

	gInpPath, _ = gtk.EntryNew()
	gInpPath.SetText(path.GetVisual())
	gInpPath.SetHExpand(true)
	gInpPath.SetHAlign(gtk.ALIGN_FILL)

	gBtnRefresh, _ = gtk.ButtonNewWithLabel("Reload")
	//img2 := GTK_Image_From_File(appdir+"gui/button_reload.png", "png")
	img2 := GTK_Image_From_Name("view-refresh", gtk.ICON_SIZE_BUTTON)
	gBtnRefresh.SetImage(img2)
	gBtnRefresh.SetProperty("always-show-image", true)
	gBtnRefresh.Connect("clicked", func() {
		tpath, _ := gInpPath.GetText()
		path.SetVisual(tpath)
		listFiles(gGFiles, path.GetReal())
	})

	// gGTop, _ := gtk.GridNew()
	// gGTop.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	// gGTop.Attach(gBtnUp, 0, 0, 1, 1)
	// gGTop.Attach(gInpPath, 1, 0, 1, 1)
	// gGTop.Attach(gBtnRefresh, 2, 0, 1, 1)

	header, _ := gtk.HeaderBarNew()
	//header.Add(gBtnUp)
	header.PackStart(gBtnUp)
	//header.Add(gInpPath)
	header.SetCustomTitle(gInpPath)
	//header.Add(gBtnRefresh)
	header.PackEnd(gBtnRefresh)
	header.SetHExpand(true)

	// ================

	//gGDiscs, _ = gtk.GridNew()
	gGDiscs, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	gGDiscs.SetOrientation(gtk.ORIENTATION_VERTICAL)
	gGDiscs.SetHExpand(true)
	//gGDiscs.SetSizeRequest(LEFT_PANEL_SIZE, 0)

	sLeftScroll, _ = gtk.ScrolledWindowNew(nil, nil)
	sLeftScroll.Add(gGDiscs)
	sLeftScroll.SetVExpand(true)
	sLeftScroll.SetHExpand(true)
	//sLeftScroll.SetSizeRequest(LEFT_PANEL_SIZE, 0)
	sLeftScroll.Connect("draw", func(g *gtk.ScrolledWindow, ctx *cairo.Context) {
		ctx.SetSourceRGBA(0.5, 0.5, 0.5, 0.5)
		//g.CheckResize()
		h := g.GetAllocatedHeight()
		//Prln("h" + I2S(h))
		lw := sLeftScroll.GetAllocatedWidth()
		//lw=LEFT_PANEL_SIZE
		ctx.Rectangle(0, 0, float64(lw), float64(h-2))
		ctx.Fill()
	})

	gGFiles, _ = gtk.GridNew()
	gGFiles.SetOrientation(gtk.ORIENTATION_VERTICAL)
	gGFiles.SetHAlign(gtk.ALIGN_START)
	//gGFiles.SetVExpand(true)
	gGFiles.SetHExpand(true)
	gGFiles.SetColumnHomogeneous(true)

	//hadjusment:=gtk.AdjustmentNew()
	sRightScroll, _ = gtk.ScrolledWindowNew(nil, nil)
	sRightScroll.SetVExpand(true)
	sRightScroll.SetHExpand(true)
	sRightScroll.Add(gGFiles)

	rightEv, _ := gtk.EventBoxNew()
	rightEv.Connect("button-press-event", func(_ *gtk.ScrolledWindow, event *gdk.Event) {
		eventbutton := &gdk.EventButton{event}
		mousekey := eventbutton.ButtonVal()
		if mousekey == 3 {
			if rightmenu == nil || !rightmenu.IsVisible() {
				rightmenu, _ = gtk.MenuNew()
				Menu_CurrentFolder(rightmenu, path.GetReal())
				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event)
			} else {
				Prln("ignoring menu")
			}
		}
	})
	rightEv.Add(sRightScroll)

	//gGCenter, _ := gtk.GridNew()
	//gGCenter.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	//gGCenter.Attach(sLeft, 0, 1, 1, 1)
	//gGCenter.Attach(sScroll, 1, 1, 1, 1)

	gGCenter, _ = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	gGCenter.Add(sLeftScroll)
	gGCenter.Add(rightEv)
	gGCenter.SetHExpand(true)
	gGCenter.SetPosition(LEFT_PANEL_SIZE)
	ps := 0
	lw_old := 0
	gGCenter.Connect("size-allocate", func() {
		ps += 1
		lw := gGCenter.GetPosition()
		if lw_old != lw {
			lw_old = lw
			LEFT_PANEL_SIZE = lw - 5
			resize_files_icons()
		}
	})

	// ================

	rezoom := func() {
		GTK_Childs(gGFiles, true, true)
		//path, _ = gInpPath.GetText()
		listFiles(gGFiles, path.GetReal())
		resize_files_icons()
	}

	/*gBtnZOOM, _ := gtk.ButtonNewWithLabel("zoom")
	gBtnZOOM.Connect("clicked", func() {
		//ZOOM_SIZE = 192 - ZOOM_SIZE
		ZOOM_SIZE = (ZOOM_SIZE * 2) % (512 - 64)
		Prln("zoom:" + I2S(ZOOM_SIZE))
		rezoom()
	})*/

	spin, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 1, 3, 1)
	spin.SetSizeRequest(90, 30)
	spin.SetDrawValue(false)
	//sp.SetSensitive(false)
	switch ZOOM_SIZE {
	case 64:
		spin.SetValue(1)
	case 128:
		spin.SetValue(2)
	case 256:
		spin.SetValue(3)
	}
	spin.Connect("value-changed", func() {
		sv := RoundF(spin.GetValue())
		old_zoom := ZOOM_SIZE
		switch sv {
		case 1:
			ZOOM_SIZE = 64
		case 2:
			ZOOM_SIZE = 128
		case 3:
			ZOOM_SIZE = 256
		}
		//sp.SetDrawValue("x" + I2S(ZOOM_SIZE))
		if old_zoom != ZOOM_SIZE {
			Prln("ss" + I2S(int(sv)))
			rezoom()
			spin.SetValue(float64(sv))
		}
	})
	//sp.SetDrawValue(false)

	mem, _ = gtk.LabelNew("MEM")
	space, _ = gtk.LabelNew("")
	space.SetHExpand(true)
	GTK_LabelWrapMode(space, 2)

	/*gBtnGarbage, _ := gtk.ButtonNewWithLabel("GarbageCollection")
	gBtnGarbage.Connect("clicked", func() {
		GarbageCollection()
		FreeOSMemory()
		Prln("FreeOSMemory()")
	})*/

	gCheckDragCopy, _ := gtk.CheckButtonNewWithLabel("mouse drag works as copy")
	gCheckDragCopy.SetActive(false)
	gCheckDragCopy.SetSensitive(false)
	gCheckDragCopy.Connect("clicked", func() {
		//copy_mode = gCheckDragCopy.GetActive()

	})

	gCheckPreviewFolders, _ := gtk.CheckButtonNewWithLabel("preview folders")
	gCheckPreviewFolders.SetActive(with_folders_preview)
	gCheckPreviewFolders.Connect("clicked", func() {
		with_folders_preview = gCheckPreviewFolders.GetActive()
		listFiles(gGFiles, path.GetReal())
	})
	gCheckPreviewFiles, _ := gtk.CheckButtonNewWithLabel("preview files")
	gCheckPreviewFiles.SetActive(with_files_preview)
	gCheckPreviewFiles.Connect("clicked", func() {
		with_files_preview = gCheckPreviewFiles.GetActive()
		listFiles(gGFiles, path.GetReal())
	})

	gGDown, _ := gtk.GridNew()
	gGDown.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	//gGDown.Attach(gBtnGarbage, 0, 0, 1, 1)
	gGDown.Attach(mem, 1, 0, 1, 1)
	gGDown.Attach(space, 2, 0, 1, 1)
	gGDown.Attach(gCheckDragCopy, 3, 0, 1, 1)
	gGDown.Attach(gCheckPreviewFolders, 4, 0, 1, 1)
	gGDown.Attach(gCheckPreviewFiles, 5, 0, 1, 1)
	gGDown.Attach(spin, 6, 0, 1, 1)

	// =================

	GTK_CopyPasteDnd_SetAppDest(sRightScroll)

	//Prln(I2S(int(gdk.KEY_c)))

	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		GTK_CopyPasteDnd_SetWindowKeyPressed(path, ev)
	})
	// win.Connect("key-release-event", func(win *gtk.Window, ev *gdk.Event) {
	// 	keyEvent := &gdk.EventKey{ev}
	// 	uint_key := keyEvent.KeyVal()
	// 	Prln("keyup:" + I2S(int(uint_key)))
	// })

	// ================

	menuBar := GTK_MainMenu(win)

	// ================

	gGMain, _ := gtk.GridNew()
	gGMain.SetOrientation(gtk.ORIENTATION_VERTICAL)
	gGMain.Attach(menuBar, 0, 0, 1, 1)
	gGMain.Attach(header, 0, 1, 1, 1)
	//win.SetTitlebar(header)
	gGMain.Attach(gGCenter, 0, 2, 1, 1)
	gGMain.Attach(gGDown, 0, 3, 1, 1)

	win.Add(gGMain)
	win.ShowAll()

	num_threads := opt.GetThreads()
	//RuntimeGoMaxProcs(num_threads)
	go IconThread(icon_chan1, qu, 0)
	for t := 1; t <= num_threads; t++ {
		go IconThread(icon_chanN, qu, t)
	}

	listDiscs(gGDiscs)
	listFiles(gGFiles, path.GetReal())

	pid := AppProcessID()
	Prln("PID:" + I2S(pid))

	go func() {
		for {
			SleepMS(1500)
			GarbageCollection()
			FreeOSMemory()
			usage = F2S(LinixMemoryUsedMB(pid), 1) + "Mb"
		}
	}()

	fswatcher.SetListenerOnce()

	MainThread()
	//gtk.Main()

	/*for {
		Sleep(500)
		//fileMutex.Lock()
		//mem.SetText("RAM Usage: " + I2S(GetPC_MemoryUsageMb()) + " Mb")
		//fileMutex.Unlock()
	}*/
	//select {}

}

func listDiscs(g *gtk.Box) {

	GTK_Childs(g, true, true) //arrd :=
	//Prln("disc_child_len:" + I2S(len(arrd)))

	discs := GetDiscParts(true, true, true, true)

	mountlist = LinuxGetMountList()

	for j, di := range discs {
		d := di
		Prln(d.String())
		p := d.MountPath
		gDBtn, _ := gtk.ButtonNew() //WithLabel(d.Title + "\n" + d.Model + "\n" + d.PartName)
		gDBtn.SetHExpand(true)
		//gDBtn.SetSizeRequest(LEFT_PANEL_SIZE, 32)
		//gDBtn.SetHAlign(gtk.ALIGN_CENTER)
		//gDBtn.SetJustify(gtk.JUSTIFY_CENTER)
		gDBtn.Connect("clicked", func() {
			//tpath := NewLinuxPath(true)
			path.SetReal(p)
			gInpPath.SetText(path.GetVisual())
			listFiles(gGFiles, path.GetReal())
		})
		gDBtn.SetHExpand(false)

		grid, _ := gtk.GridNew()
		grid.SetHExpand(true)

		lbl1, _ := gtk.LabelNew(d.Title)
		lbl1.SetJustify(gtk.JUSTIFY_LEFT)
		lbl1.SetHAlign(gtk.ALIGN_START)
		GTK_LabelWrapMode(lbl1, 1)
		lbl1.SetLineWrap(true)
		grid.Attach(lbl1, 0, 0, 1, 1)
		lbl1.SetMarkup("<b><u>" + HtmlEscape(d.Title) + "</u></b>")

		lbl2, _ := gtk.LabelNew(d.SpaceTotal)
		lbl2.SetJustify(gtk.JUSTIFY_RIGHT)
		lbl2.SetHAlign(gtk.ALIGN_END)
		//lbl2.SetHExpand(true)
		grid.Attach(lbl2, 1, 0, 1, 1)

		if j == 0 {
			lbl1.SetMarkup("<b><u>" + HtmlEscape(StringUp(d.PartName)) + "</u></b>")
			lbl2.SetText("HOME")
		}

		extra := StringFind(StringDown(d.PartName), "gvfsd-fuse") > 0
		if len(StringTrim(d.Model)) > 0 && !extra {
			lbl3, _ := gtk.LabelNew(d.Model)
			lbl3.SetJustify(gtk.JUSTIFY_LEFT)
			lbl3.SetHAlign(gtk.ALIGN_START)
			lbl3.SetHExpand(true)
			GTK_LabelWrapMode(lbl3, 1)
			grid.Attach(lbl3, 0, 1, 2, 1)
		}

		if j > 0 && d.SpacePercent > -1 {
			lbl4, _ := gtk.LabelNew(d.PartName)
			if extra {
				lbl4.SetText(d.Model)
			}
			lbl4.SetJustify(gtk.JUSTIFY_LEFT)
			lbl4.SetHAlign(gtk.ALIGN_START)
			lbl4.SetHExpand(true)
			GTK_LabelWrapMode(lbl4, 1)
			grid.Attach(lbl4, 0, 2, 1, 1)

			fs := d.Protocol
			if d.Protocol == "PART" {
				fs = StringUp(d.FSType)
			}
			lbl5, _ := gtk.LabelNew(fs)
			lbl5.SetMarkup("<u>" + HtmlEscape(fs) + "</u>")
			lbl5.SetJustify(gtk.JUSTIFY_RIGHT)
			lbl5.SetHAlign(gtk.ALIGN_END)
			//lbl4.SetHExpand(true)
			grid.Attach(lbl5, 1, 2, 1, 1)

			progr, _ := gtk.LevelBarNew() // ProgressBarNew()
			//progr.SetFraction(float64(d.SpacePercent) / 100.0)
			progr.SetValue(float64(d.SpacePercent) / 100.0)
			progr.SetHExpand(true)
			progr.SetSizeRequest(10, 0)

			grid.Attach(progr, 0, 3, 2, 1)
		}

		gDBtn.Connect("button-press-event", func(_ *gtk.Button, event *gdk.Event) {
			eventbutton := &gdk.EventButton{event}
			mousekey := eventbutton.ButtonVal()
			switch mousekey {
			case 3:
				Prln("right")
				rightmenu, _ := gtk.MenuNew()

				if !d.Primary && d.SpacePercent >= 0 {
					GTK_MenuItem(rightmenu, "Eject", nil)
				}
				if d.SpacePercent < 0 {
					GTK_MenuItem(rightmenu, "Remove", nil)
				}
				GTK_MenuItem(rightmenu, "Info", nil)

				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event) // (evBox, gdk.GDK_GRAVITY_STATIC, gdk.GDK_GRAVITY_STATIC,
			}
		})

		gDBtn.Add(grid)
		//g.Attach(gDBtn, 0, j, 1, 1)
		g.Add(gDBtn)
	}

	g.ShowAll()
}

func listFiles(g *gtk.Grid, lpath string) {

	//sRightScroll.

	fswatcher.Select(lpath)

	new_ind := path_updated.Add(1)

	if with_destroy {
		for j := 0; j < len(arr_blocks); j++ {
			arr_blocks[j].Destroy()
			arr_blocks[j] = nil
		}
	}

	arr_blocks = []*GtkFileIconBlock{}

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

	files, err := Folder_ListFiles(lpath2)
	if err != nil {
		Prln(err.Error())
		iconwithlabel := NewFileIconBlock(lpath2, "ERROR!", 400, false, false, false, false, err.Error())
		arr_blocks = append(arr_blocks, iconwithlabel)
		g.Attach(iconwithlabel.GetWidget(), 1, 1, 1, 1)
		g.ShowAll()
		return
	}
	j := 0

	SortArray(files, func(i, j int) bool {
		if files[i].IsDir() != files[j].IsDir() {
			return !CompareBoolLess(files[i].IsDir(), files[j].IsDir())
		}
		if files[i].Name() != files[j].Name() {
			return FileSortName(StringDown(files[i].Name())) < FileSortName(StringDown(files[j].Name()))
		}
		return false
	})

	var arr_render []*IconUpdateable
	icon_block_max_n, icon_block_max_w := max_icon_n_w()

	for _, f := range files {
		fname := f.Name()
		isdir := f.IsDir()
		isapp := false
		islink := FileIsLink(f)
		isregular := f.Mode().IsRegular() || islink
		oldbuf := false

		if islink {
			isdir = FileLinkIsDir(lpath2 + fname)
		}

		filepathfinal := lpath2 + fname
		if isdir {
			filepathfinal = FolderPathEndSlash(filepathfinal)
		}

		//Prln("[" + B2S_YN(isdir) + "]:{" + fname + "}" + B2S_YN(islink) + "/" + f.Mode().String())

		x := j % icon_block_max_n
		y := j / icon_block_max_n

		inf := "" //f.Mode().String() + "\n" // + "|" + f.Mode().Perm().String()

		not_read := false
		if isdir {
			if !single_thread_protocol && with_extra_info {
				fl, err := Folder_ListFiles(filepathfinal)
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
		iconwithlabel := NewFileIconBlock(lpath2, fname, icon_block_max_w, isdir, islink, not_read, ismount, inf)

		if isdir {
			if filepathfinal == opt.GetHashFolder() {
				iconwithlabel.SetIconPixPuf(GetIcon_PixBif_OF(ZOOM_SIZE, PREFIX_DRAWONME+FILE_TYPE_FOLDER_HASH))
			} else {
				hashpix := ReadHashPixbuf(filepathfinal, ZOOM_SIZE)
				if hashpix != nil {
					iconwithlabel.SetIconPixPuf(hashpix)
					oldbuf = true
				} else {
					iconwithlabel.SetIconPixPuf(GetIcon_PixBif(ZOOM_SIZE, "", true))
				}
			}
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
				if path.GetReal() != opt.GetHashFolder() {
					pixbuf_icon = ReadHashPixbuf(filepathfinal, ZOOM_SIZE)
				}
				if pixbuf_icon != nil {
					oldbuf = true
				} else {
					pixbuf_icon = GetIcon_PixBif(ZOOM_SIZE, tfile, false)
				}
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

		clicktime := TimeAddMS(TimeNow(), -2000)

		iconwithlabel.ConnectEventBox("button-press-event", func(_ *gtk.EventBox, event *gdk.Event) {
			eventbutton := &gdk.EventButton{event}
			mousekey := eventbutton.ButtonVal()
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
						listFiles(gGFiles, path.GetReal())
					} else {
						OpenFileByApp(path.GetReal()+txtlbl, "")
					}
				} else {
					clicktime = TimeNow()
				}
			case 3:
				Prln("right")
				if rightmenu != nil && rightmenu.IsVisible() {
					Prln("hiding menu")
					rightmenu.Destroy()
				}
				rightmenu, _ = gtk.MenuNew()

				Menu_FilesContextMenu(rightmenu, lpath2, fname, isdir, isapp)

				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event) // (evBox, gdk.GDK_GRAVITY_STATIC, gdk.GDK_GRAVITY_STATIC,
			}
		})

		arr_blocks = append(arr_blocks, iconwithlabel)
		g.Attach(iconwithlabel.GetWidget(), x, y, 1, 1)
		j++

		if isdir {
			fullname := FolderPathEndSlash(path.GetReal() + fname)
			if with_folders_preview && !single_thread_protocol && fullname != opt.GetHashFolder() {
				iconwithlabel.SetLoading(true)
				arr_render = append(arr_render, &IconUpdateable{icon: iconwithlabel.icon, loading: iconwithlabel.icon_loading, fullname: fullname, fname: fname, tfile: "", basic_mode: single_thread_protocol, folder: true, oldbuf: oldbuf})
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
						iconwithlabel.SetLoading(true)
						arr_render = append(arr_render, &IconUpdateable{icon: iconwithlabel.icon, loading: iconwithlabel.icon_loading, fullname: fullname, fname: fname, tfile: tfile, basic_mode: single_thread_protocol, folder: false, oldbuf: oldbuf})
					}
				}
			}
		}
	}
	g.ShowAll()
	//win.QueueDraw()

	//Prln("folder loaded. starting chans sending...")

	go func() {
		SleepMS(5)
		for b := 0; b < 2; b++ {
			oldbuf := b == 1
			for k := 0; k < 2; k++ {
				fold := k == 1
				for j := 0; j < len(arr_render); j++ {
					if arr_render[j].folder == fold && arr_render[j].oldbuf == oldbuf {
						if new_ind == path_updated.Get() {
							if single_thread_protocol {
								icon_chan1 <- arr_render[j]
							} else {
								icon_chanN <- arr_render[j]
							}
						}
						RuntimeGosched()
					}
				}
			}
		}

		/*if !single_thread_protocol {
			for j := 0; j < len(arr_render); j++ {
				if arr_render[j].folder {
					a := arr_render[j]
					arr_render2 := &IconUpdateable{icon: a.icon, fullname: a.fullname, fname: a.fname, tfile: a.tfile, basic_mode: false, folder: true}
					icon_chanN <- arr_render2
					RuntimeGosched()
				}
			}
		}*/

		//Prln("GO FINISH")
	}()

}

/*func update_icon(lpath2 string, fname string, icon *gtk.Image) {
	tfile := GetFileExtension(fname)
	if len(tfile) > 0 {
		fullname := FolderPathEndSlash(path) + fname
		iconpath := FileIconBySystem(fullname)
		if len(iconpath) > 0 {
			pixbuf_preview, err := gdk.PixbufNewFromFile(iconpath)
			if err == nil {
				pixbuf_preview2, ok := ResizePixelBuffer(pixbuf_preview, ZOOM_SIZE, gdk.INTERP_BILINEAR)
				if ok {
					icon.SetFromPixbuf(pixbuf_preview2)
				} else {
					icon.SetFromPixbuf(pixbuf_preview)
				}
			}
		}
	}
}*/

func resize_files_icons() {
	//ev++
	//Prln("resize" + I2S(ev))

	icon_block_max_n, icon_block_max_w := max_icon_n_w()
	if icon_block_max_n_old != icon_block_max_n || icon_block_max_w_old != icon_block_max_w {
		Prln("resized")
		//fileMutex.Lock()
		//defer fileMutex.Unlock()

		icon_block_max_n_old = icon_block_max_n
		icon_block_max_w_old = icon_block_max_w
		icon_block_max_w += 0

		arr := GTK_Childs(gGFiles, true, false)
		len_arr := len(arr)
		for j := 0; j < len_arr; j++ {
			x := j % icon_block_max_n
			y := j / icon_block_max_n
			gEv := arr[len_arr-j-1]
			gGFiles.Attach(gEv, x, y, 1, 1)

			/*label, ok := map_labels[&gEv]
			if ok {
				label.SetSizeRequest(icon_block_max_w-BORDER_SIZE*2, 32)
				txt, _ := label.GetText()
				Prln(txt)
			}*/
			/*Prln(Typeof(gEv))
			type Resizer interface {
				SetSizeRequest(a int, b int)
			}
			if r, ok := gEv.(Resizer); ok {
				Prln("ok_" + I2S(j))
				r.SetSizeRequest(icon_block_max_w-BORDER_SIZE*2, 32)
			} else {
				Prln("fail_" + I2S(j))
			}*/

		}

		len_arrl := len(arr_blocks)
		for j := 0; j < len_arrl; j++ {
			arr_blocks[j].SetWidth(icon_block_max_w)
		}
		//gGFiles.ShowAll()
	}
}

func max_icon_n_w() (int, int) {
	//ww, _ := sScroll.GetPreferredWidth()
	//sScroll.CheckResize()
	//ww, _ := win.GetPreferredWidth()
	ww, _ := win.GetSize()
	real_w := MAXI(16, ww-LEFT_PANEL_SIZE) - 6
	icon_block_max_w := MAXI(16, ZOOM_SIZE+BORDER_SIZE*4)
	icon_block_max_n := MAXI(1, MAXI(16, real_w)/icon_block_max_w)
	icon_block_max_w = real_w/icon_block_max_n - BORDER_SIZE*3
	//Prln("size" + I2S(ww))
	return icon_block_max_n, icon_block_max_w
}

func MainThread() {
	iter := 0
	gtk.MainIteration()
	RuntimeGosched()
	for {
		if fswatcher.IsUpdated() {
			listFiles(gGFiles, path.GetReal())
		}
		gtk.MainIteration()
		qlen := qu.Length()
		if qlen > 0 {
			//Prln("qlen:" + I2S(qlen) + " / " + F2S(GetPC_MemoryUsageMb(), 1) + "Mb")
			//Prln("it1")
			w, ok := qu.GetEnd().(*IconUpdateable)
			for ok && !GTK_WidgetExist(w.icon) && qu.Length() > 0 {
				w, ok = qu.GetEnd().(*IconUpdateable)
				Prln("widget searching...")
			}
			if ok && GTK_WidgetExist(w.icon) {
				//Prln("pixbufset")
				w.loading.SetFromPixbuf(nil)
				w.icon.SetFromPixbuf(w.pixbuf_preview)
			}
			//Prln("it2")
		} else {
			iter++
		}
		//if iter > 10 {
		//	iter = 0
		mem.SetText(I2S(num_works.Get()) + " processes; RAM Usage: " + F2S(GetPC_MemoryUsageMb(), 1) + " Mb & " + usage)
		main_iterations_funcs.ExecAll()
		//}
		RuntimeGosched()
		//debug.FreeOSMemory()
		//mem.SetText("RAM Usage: " + I2S(linux.LinuxMemory()) + " Mb")
		//GarbageCollection()
		//win.ShowAll()
	}
}
