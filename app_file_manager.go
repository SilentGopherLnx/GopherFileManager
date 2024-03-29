package main

//sudo apt-get install libgtk-3-dev
//sudo apt-get install libcairo2-dev
//sudo apt-get install libglib2.0-dev

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_fileicon"
	. "github.com/SilentGopherLnx/GopherFileManager/pkg_filetools"

	//	"os/exec"

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
var gInpPath, gInpSearch *gtk.Entry
var gBtnUp *gtk.Button
var gBtnRefresh *gtk.Button
var mem, space *gtk.Label
var top_current_menu *gtk.Menu

var gBtnBack, gBtnForward *gtk.Button
var hist *PathHistory = PathHistoryNew()

var path *LinuxPath = NewLinuxPath(true)

//var path_updated = NewAtomicInt(0)
var req_id = NewAtomicInt64(0)

var icon_block_max_n_old, icon_block_max_w_old int

var ZOOM_SIZE = 64

var LEFT_PANEL_SIZE = 200 //200

var qu *SyncQueue

var with_folders_preview bool = false
var with_files_preview bool = false
var with_cache_preview bool = false

var with_destroy bool = true

//var with_extra bool = false

var usage = ""

var rightmenu *gtk.Menu = nil

var main_iterations_funcs *FuncArr = NewFuncArr()

var mountlist [][2]string = [][2]string{}

var fswatcher *FSWatcher

var num_works *AInt = NewAtomicInt(0)

var sort_reverse bool = false
var sort_mode int = 0

var sudo string = ""

var upd_func func()

//var killchan chan *exec.Cmd

var spinnerIcons *gtk.Spinner
var spinnerFiles *gtk.Spinner

var langs *LangArr

var drag_mode = false

func init() {

	RuntimeLockOSThread()
	AboutVersion(AppVersion())
	langs = InitLang(FolderPathEndSlash(FolderLocation_App()+"localization") + "translation_manager.cfg")
	InitOptions()
	ZOOM_SIZE = opt.GetZoom()

	args := AppRunArgs()
	if len(args) >= 2 {
		path.SetVisual(args[1])
	} else {
		path.SetReal(FolderLocation_WorkDir())
	}

	icon_chanN = make(chan *IconUpdateable)
	icon_chan1 = make(chan *IconUpdateable)
	qu = NewSyncQueue()
	//	killchan = make(chan *exec.Cmd)

	with_cache_preview = true
	with_folders_preview = true
	with_files_preview = true
	//with_destroy = false
}

func main() {

	fswatcher = NewFSWatcher(opt.GetInotifyPeriod())
	defer fswatcher.Close()

	gtk.Init(nil)
	Prln("gtk_version:" + GTK_GetVersion())

	var err error
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		Prln("Unable to create window:") // + err)
	}

	GTK_ColorsLoad(win)

	uid, _, _ := GetPC_UserUidLoginName()
	sudo = B2S(LinuxRoot_Check() == 1, "[root"+uid+"] ", "-")
	upd_title()
	win.SetDefaultSize(1200, 800)
	win.SetPosition(gtk.WIN_POS_CENTER)

	if AppHasArg("-max") {
		win.Maximize()
	}

	win.Connect("destroy", func() {
		AppExit(0)
	})

	//ev := 0 // https://developer.gnome.org/gtk3/stable/GtkWidget.html
	win.Connect("size-allocate", func() {
		resize_event_no_repeats()
	})

	appdir := FolderLocation_App()
	win.SetIconFromFile(appdir + GUI_PATH + "icon.png")

	// ================

	spinnerFiles, _ = gtk.SpinnerNew()

	gInpPath, _ = gtk.EntryNew()
	gInpPath.SetText(path.GetVisual())
	gInpPath.SetHExpand(true)
	gInpPath.SetHAlign(gtk.ALIGN_FILL)
	gInpPath.Connect("button-press-event", func() {
		gInpPath.SetCanFocus(true)
	})

	gInpSearch, _ = gtk.EntryNew()
	gInpSearch.SetText("")
	gInpSearch.SetHExpand(true)
	//gInpSearch.SetHAlign(gtk.ALIGN_FILL)
	gInpSearch.SetPlaceholderText(langs.GetStr("main_search") + ":")
	gInpSearch.Connect("button-press-event", func() {
		gInpSearch.SetCanFocus(true)
	})

	gBtnUp, _ = gtk.ButtonNewWithLabel(langs.GetStr("main_up"))
	//gBtnUp.SetProperty("background-color", "red")
	//img1 := GTK_Image_From_File(appdir+"gui/button_up.png", "png")
	img1 := GTK_Image_From_Name("go-up", gtk.ICON_SIZE_BUTTON)
	gBtnUp.SetImage(img1)
	gBtnUp.SetProperty("always-show-image", true)
	gBtnUp.Connect("clicked", func() {
		path.GoUp()
		gInpPath.SetText(path.GetVisual())
		gInpSearch.SetText("")
		listFiles(gGFiles, path, true, true)
	})
	gBtnUp.SetCanFocus(false)

	upd_func = func() {
		tpath, _ := gInpPath.GetText()
		path.SetVisual(tpath)
		listFiles(gGFiles, path, true, true)
	}

	gBtnBack, _ = gtk.ButtonNewWithLabel(langs.GetStr("main_back"))
	img_bk := GTK_Image_From_Name("go-previous", gtk.ICON_SIZE_BUTTON)
	gBtnBack.SetImage(img_bk)
	gBtnBack.SetProperty("always-show-image", true)
	gBtnBack.Connect("clicked", func() {
		ok, v, s := hist.Back()
		if ok {
			path.SetVisual(v)
			gInpPath.SetText(v)
			gInpSearch.SetText(s)
			listFiles(gGFiles, path, true, false)
		}
	})
	gBtnBack.SetCanFocus(false)

	gBtnForward, _ = gtk.ButtonNewWithLabel(langs.GetStr("main_forward"))
	img_fw := GTK_Image_From_Name("go-next", gtk.ICON_SIZE_BUTTON)
	gBtnForward.SetImage(img_fw)
	gBtnForward.SetProperty("always-show-image", true)
	gBtnForward.Connect("clicked", func() {
		ok, v, s := hist.Forward()
		if ok {
			path.SetVisual(v)
			gInpPath.SetText(v)
			gInpSearch.SetText(s)
			listFiles(gGFiles, path, true, false)
		}
	})
	gBtnForward.SetCanFocus(false)

	gBtnRefresh, _ = gtk.ButtonNewWithLabel(langs.GetStr("main_reload"))
	//img2 := GTK_Image_From_File(appdir+"gui/button_reload.png", "png")
	img2 := GTK_Image_From_Name("view-refresh", gtk.ICON_SIZE_BUTTON)
	gBtnRefresh.SetImage(img2)
	gBtnRefresh.SetProperty("always-show-image", true)
	gBtnRefresh.Connect("clicked", func() {
		upd_func()
	})
	gBtnRefresh.SetCanFocus(false)

	gGTop1, _ := gtk.GridNew()
	gGTop1.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	gGTop1.Attach(gBtnBack, 0, 0, 1, 1)
	gGTop1.Attach(gBtnForward, 1, 0, 1, 1)
	gGTop1.Attach(gBtnUp, 2, 0, 1, 1)
	//gGTop1.Attach(gInpPath, 3, 0, 1, 1)

	gGTop2, _ := gtk.GridNew()
	gGTop2.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	gGTop2.Attach(gInpPath, 0, 0, 1, 1)
	gGTop2.Attach(gInpSearch, 1, 0, 1, 1)
	gGTop2.Attach(spinnerFiles, 2, 0, 1, 1)
	gGTop2.Attach(gBtnRefresh, 3, 0, 1, 1)

	header, _ := gtk.HeaderBarNew()
	header.Add(gGTop1)
	//header.PackStart(gGTop1)
	//header.SetCustomTitle(gInpPath)
	header.SetCustomTitle(gGTop2)
	//header.PackEnd(gGTop2)
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
	// gGFiles.SetMarginStart(BORDER_SIZE)
	// gGFiles.SetMarginEnd(BORDER_SIZE * 3 / 2)
	// gGFiles.SetMarginTop(BORDER_SIZE)
	// gGFiles.SetMarginBottom(BORDER_SIZE)
	gGFiles.SetMarginEnd(BORDER_SIZE / 2)
	gGFiles.SetBorderWidth(uint(BORDER_SIZE))
	gGFiles.SetColumnSpacing(uint(BORDER_SIZE))
	gGFiles.SetRowSpacing(uint(BORDER_SIZE))

	//hadjusment:=gtk.AdjustmentNew()
	sRightScroll, _ = gtk.ScrolledWindowNew(nil, nil) // g_signal_connect(G_OBJECT(browserscrolledview),"scroll-event",G_CALLBACK(userActive),NULL);
	sRightScroll.SetVExpand(true)
	sRightScroll.SetHExpand(true)
	sRightScroll.Add(gGFiles)

	rightEv, _ := gtk.EventBoxNew()
	rightEv.Connect("draw", func(_ *gtk.ScrolledWindow, ctx *cairo.Context) {
		_, dy := GTK_ScrollGetValues(sRightScroll)
		FilesSelector_Draw(dy, ctx)
	})
	rightEv.Connect("button-press-event", func(_ *gtk.EventBox, event *gdk.Event) {
		disable_focus()
		mousekey, _, _, zone := FilesSelector_MousePressed(event, sRightScroll)
		if mousekey == 3 && zone {
			if rightmenu == nil || !rightmenu.IsVisible() {
				rightmenu, _ = gtk.MenuNew()
				_, s := hist.GetCurrent()
				GTKMenu_CurrentFolder(rightmenu, *path, s)
				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event)
			} else {
				Prln("ignoring menu")
			}
		}
	})
	rightEv.Connect("motion-notify-event", func(_ *gtk.EventBox, event *gdk.Event) {
		FilesSelector_MouseMoved(event, sRightScroll, win.QueueDraw)
	})
	rightEv.Connect("button-release-event", func(_ *gtk.EventBox, event *gdk.Event) {
		FilesSelector_MouseRelease(event, sRightScroll, win.QueueDraw)
	})
	sRightScroll.SetEvents(int(gdk.ALL_EVENTS_MASK))
	/*	sn := 0
		sRightScroll.Connect("scroll-event", func(_ *gtk.ScrolledWindow, event *gdk.Event) {
			sn++
			Prln("scroll -" + I2S(sn))
			if select_x1 > 0 && select_y1 > 0 {
				eventscroll := &gdk.EventScroll{event}
				Prln("scroll " + F2S(eventscroll.DeltaY(), 2))
				//win.QueueDraw()
			}
		})*/
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
			resize_event_no_repeats()
		}
	})

	// ================

	rezoom := func() {
		GTK_Childs(gGFiles, true, true)
		//path, _ = gInpPath.GetText()
		listFiles(gGFiles, path, true, false)
		resize_event_no_repeats()
	}

	spin_title, _ := gtk.LabelNew(langs.GetStr("gui_down_zoom") + ":")
	spin_title.SetMarkup("<b>" + langs.GetStr("gui_down_zoom") + ":</b>")
	za := Constant_ZoomArray()
	spin, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 1, float64(len(za)), 1)
	spin.SetSizeRequest(90, 30)
	spin.SetDrawValue(false)
	//sp.SetSensitive(false)
	zi := IntInArray(ZOOM_SIZE, za)
	if zi > -1 {
		spin.SetValue(float64(zi) + 1)
	} else {
		spin.SetValue(1)
	}
	spin.Connect("value-changed", func() {
		old_zoom := ZOOM_SIZE
		sv := RoundF(spin.GetValue())
		if sv > 0 && sv < len(za) {
			ZOOM_SIZE = za[sv-1]
		} else {
			if sv < 1 {
				ZOOM_SIZE = za[0]
			} else {
				ZOOM_SIZE = za[len(za)-1]
			}
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

	gCheckPreviewCache, _ := gtk.CheckButtonNewWithLabel("preview cache")
	gCheckPreviewCache.SetActive(with_cache_preview)
	gCheckPreviewCache.Connect("clicked", func() {
		with_cache_preview = gCheckPreviewCache.GetActive()
		if with_cache_preview {
			listFiles(gGFiles, path, true, false)
		}
	})

	gCheckPreviewFolders, _ := gtk.CheckButtonNewWithLabel("preview folders")
	gCheckPreviewFolders.SetActive(with_folders_preview)
	gCheckPreviewFolders.Connect("clicked", func() {
		with_folders_preview = gCheckPreviewFolders.GetActive()
		if with_folders_preview {
			listFiles(gGFiles, path, true, false)
		}
	})
	gCheckPreviewFiles, _ := gtk.CheckButtonNewWithLabel("preview files")
	gCheckPreviewFiles.SetActive(with_files_preview)
	gCheckPreviewFiles.Connect("clicked", func() {
		with_files_preview = gCheckPreviewFiles.GetActive()
		if with_files_preview {
			listFiles(gGFiles, path, true, false)
		}
	})

	spinnerIcons, _ = gtk.SpinnerNew()

	gGDown, _ := gtk.GridNew()
	gGDown.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	//gGDown.Attach(gBtnGarbage, 0, 0, 1, 1)
	gGDown.Attach(spinnerIcons, 0, 0, 1, 1)
	gGDown.Attach(mem, 1, 0, 1, 1)
	gGDown.Attach(space, 2, 0, 1, 1)
	//gGDown.Attach(gCheckDragCopy, 3, 0, 1, 1)
	gGDown.Attach(gCheckPreviewCache, 4, 0, 1, 1)
	gGDown.Attach(gCheckPreviewFolders, 5, 0, 1, 1)
	gGDown.Attach(gCheckPreviewFiles, 6, 0, 1, 1)
	gGDown.Attach(spin_title, 7, 0, 1, 1)
	gGDown.Attach(spin, 8, 0, 1, 1)

	// =================

	GTK_CopyPasteDnd_SetAppDest(&sRightScroll.Widget)

	//Prln(I2S(int(gdk.KEY_c)))

	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		if !gInpPath.IsFocus() && !gInpSearch.IsFocus() { //gGFiles.HasVisibleFocus() || gGFiles.HasFocus() || gGFiles.IsFocus() {
			key, state, hkey := GTK_KeyboardKeyOfEvent(ev)
			GTK_CopyPasteDnd_SetWindowKeyPressed(path, key, state, hkey)
		}
	})
	// win.Connect("key-release-event", func(win *gtk.Window, ev *gdk.Event) {
	// 	keyEvent := &gdk.EventKey{ev}
	// 	uint_key := keyEvent.KeyVal()
	// 	Prln("keyup:" + I2S(int(uint_key)))
	// })

	// ================

	var menuBar *gtk.MenuBar
	menuBar, top_current_menu = GTKMenu_Main(win)

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
	go Thread_Icon(icon_chan1, qu, 0)
	for t := 1; t <= num_threads; t++ {
		go Thread_Icon(icon_chanN, qu, t)
	}

	listDiscs(gGDiscs)
	listFiles(gGFiles, path, true, true)

	pid := AppProcessID()
	Prln("PID:" + I2S(pid))

	go Thread_GC_and_Free(pid)

	fswatcher.SetListenerOnce()

	Thread_Main()
	//gtk.Main()

}

func disable_focus() {
	gInpPath.SetCanFocus(false)
	gInpSearch.SetCanFocus(false)
	sRightScroll.GrabFocus()
}

func upd_title() {
	folder_name := path.GetLastNode()
	url := path.GetUrl()
	is_smb, pc_name, netfolder, smb_hasmore := SMB_CheckPath(url)
	if is_smb {
		if StringLength(pc_name) == 0 {
			folder_name = "SMB://"
		}
		if StringLength(pc_name) > 0 && StringLength(netfolder) == 0 {
			folder_name = pc_name
		}
		if StringLength(netfolder) > 0 && !smb_hasmore {
			folder_name = netfolder
		}
	}
	search := ""
	s, _ := gInpSearch.GetText()
	if s != "" {
		search = langs.GetStr("main_search") + ": [" + s + "] /"
	}
	win.SetTitle(search + folder_name + " " + sudo + " GopherFileManager")
	if top_current_menu != nil && path != nil {
		GTKMenu_CurrentFolder(top_current_menu, *path, s)
		top_current_menu.ShowAll()
	}
}

func OpenManager(path_to_folder string) {
	go ExecCommandBash("" + ExecQuote(AppRunArgs()[0]) + " " + ExecQuote(FolderPathEndSlash(path_to_folder)) + B2S(win.IsMaximized(), " -max", ""))
}

func GetLastVerison(url string) (bool, string) {
	txt, err := NetReadUrlText(url)
	//txt, _ = FileTextRead(".........../src/github.com/SilentGopherLnx/GopherFileManager/version.go")
	//txt = "\n\nconst app_version_manager = \"0.1.14\"\n\n"
	if err != nil {
		return false, err.Error()
	}
	arr := StringSplitLines(txt)
	for j := 0; j < len(arr); j++ {
		if StringFind(arr[j], "const app_version_") == 1 {
			vstr := StringSplit(arr[j], "\"")
			if len(vstr) >= 2 {
				return true, vstr[1]
			}
		}
	}
	return false, "No version file on url"
}

func AppVersion_Mover(path string) string {
	a, b, c := ExecCommand(path, "-v")
	if StringLength(b) == 0 && StringLength(c) == 0 {
		v := StringSplitLines(a)
		if len(v) > 0 {
			return v[0]
		} else {
			return "Not installed?"
		}
	} else {
		return "Not installed? " + a + b + c
	}
}
