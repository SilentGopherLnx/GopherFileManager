package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const OPER_COPY = "copy"
const OPER_MOVE = "move"
const OPER_DELETE = "delete"

// https://github.com/geany/geany/issues/1368
// GtkAccelGroup *accel_group = gtk_accel_group_new ();
// gtk_window_add_accel_group (GTK_WINDOW (window), accel_group);
// gtk_accel_group_connect (accel_group, GDK_KEY_Q, GDK_CONTROL_MASK, 0, g_cclosure_new_swap (G_CALLBACK (hello), window, NULL));
func GTK_CopyPasteDnd_SetWindowKeyPressed(path *LinuxPath, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	uint_key := keyEvent.KeyVal()
	//Prln("key:" + I2S(int(uint_key)))
	if keyEvent.State() == gdk.GDK_CONTROL_MASK { // //key:65507 Ctrl
		if uint_key == gdk.KEY_x { //120
			Prln("Ctrl+X")
			GTK_CopyPasteDnd_CopyDel(path.GetReal(), true, false)
		}
		if uint_key == gdk.KEY_c { //99
			Prln("Ctrl+C")
			GTK_CopyPasteDnd_CopyDel(path.GetReal(), false, false)
		}
		if uint_key == gdk.KEY_v { //118
			Prln("Ctrl+V")
			GTK_CopyPasteDnd_Paste(path.GetReal())
		}
		Prln(I2S(int(uint_key)))
	} else {
		if uint_key == gdk.KEY_F2 { //65471
			Prln("F2")
		}
		if uint_key == gdk.KEY_Delete { //65535
			Prln("Del")
			GTK_CopyPasteDnd_CopyDel(path.GetReal(), false, true)
		}
	}
}

func GTK_CopyPasteDnd_CopyDel(folderpath string, cut_mode bool, del bool) {
	fnames := FileSelector_GetList()
	fpath2 := FolderPathEndSlash(folderpath)
	list := []*LinuxPath{}
	for j := 0; j < len(fnames); j++ {
		file1 := NewLinuxPath(false) //??
		file1.SetReal(fpath2 + fnames[j])
		list = append(list, file1)
	}
	if !del {
		LinuxClipBoard_CopyFiles(list, cut_mode)
	} else {
		RunFileOperaion(list, nil, OPER_DELETE)
	}
}

func GTK_CopyPasteDnd_Paste(folderpath string) {
	res := LinuxClipBoard_PasteFiles()
	Prln("GTK_CopyPasteDnd_Paste:" + res)
	res_arr := StringSplitLines(res)
	if len(res_arr) > 1 {
		oper := res_arr[0]
		res_nocmd := []*LinuxPath{}
		for j := 1; j < len(res_arr); j++ {
			tpath := NewLinuxPath(false) //??
			tpath.SetUrl(res_arr[j])
			res_nocmd = append(res_nocmd, tpath)
		}
		RunFileOperaion(res_nocmd, path, StringReplace(oper, "cut", OPER_MOVE))
	}
	//res = UrlQueryUnescape(res)
	lines := StringSplitLines(res)
	if len(lines) > 3 {
		res = StringJoin(append(lines[:3], "..."), "\n")
		//space.SetText(res)
		Prln(res)
	}
}

func GTK_CopyPasteDnd_SetAppDest(w interface {
	DragDestSet(gtk.DestDefaults, []gtk.TargetEntry, gdk.DragAction)
	Connect(string, interface{}, ...interface{}) (glib.SignalHandle, error)
}) {
	t_uri, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	w.DragDestSet(gtk.DEST_DEFAULT_ALL, []gtk.TargetEntry{*t_uri}, gdk.ACTION_COPY)
	w.Connect("drag-data-received", func(g gtk.IWidget, ctx *gdk.DragContext, x int, y int, data_pointer uintptr, info uint, _ uint) {
		dnd_str := string(gtk.GetData(data_pointer))

		//space.SetText(dnd_str)
		Prln("d&d: received from another app: " + dnd_str)

		oper := OPER_MOVE
		dnd_arr := StringSplitLines(dnd_str)
		from_url := []*LinuxPath{}
		for j := 0; j < len(dnd_arr); j++ {
			tpath := NewLinuxPath(false) //??
			tpath.SetUrl(dnd_arr[j])
			from_url = append(from_url, tpath)
		}
		RunFileOperaion(from_url, path, oper)
	})
}

func GTK_CopyPasteDnd_SetIconSource(w interface {
	DragSourceSet(gdk.ModifierType, []gtk.TargetEntry, gdk.DragAction)
	Connect(string, interface{}, ...interface{}) (glib.SignalHandle, error)
}, icon *gtk.Image, getter func() []*LinuxPath) {

	t_uri_same, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_SAME_APP, 0)
	t_uri_other, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	t_uris := []gtk.TargetEntry{*t_uri_same, *t_uri_other}
	w.DragSourceSet(gdk.GDK_BUTTON1_MASK, t_uris, gdk.ACTION_COPY) //gdk.ACTION_MOVE
	w.Connect("drag-begin", func(g gtk.IWidget, ctx *gdk.DragContext) {
		Prln("d&d: drag-begin")
		pixbuf := icon.GetPixbuf()
		w := pixbuf.GetWidth()
		h := pixbuf.GetHeight()
		gtk.DragSetIconPixbuf(ctx, pixbuf, w/2, h/2)
	})
	w.Connect("drag-data-get", func(g gtk.IWidget, ctx *gdk.DragContext, data_pointer uintptr, _ uint, _ uint) {
		Prln("d&d: drag-data-get")
		files := getter()
		cmd := ""
		for j := 0; j < len(files); j++ {
			cmd = cmd + files[j].GetUrl() + "\n"
		}
		//Prln(cmd)
		gtk.SetData(data_pointer, gdk.SELECTION_TYPE_STRING, []byte(cmd))
	})
}

func GTK_CopyPasteDnd_SetFolderDest(w interface {
	DragDestSet(gtk.DestDefaults, []gtk.TargetEntry, gdk.DragAction)
	Connect(string, interface{}, ...interface{}) (glib.SignalHandle, error)
}, dest *LinuxPath) {

	t_uri_same, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_SAME_APP, 0)
	t_uri_other, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	t_uris := []gtk.TargetEntry{*t_uri_same, *t_uri_other}
	w.DragDestSet(gtk.DEST_DEFAULT_ALL, t_uris, gdk.ACTION_COPY)
	w.Connect("drag-data-received", func(g gtk.IWidget, ctx *gdk.DragContext, x int, y int, data_pointer uintptr, _ uint, _ uint) {
		dnd_str := string(gtk.GetData(data_pointer))
		oper := OPER_MOVE
		dnd_arr := StringSplitLines(dnd_str)
		from_url := []*LinuxPath{}
		for j := 0; j < len(dnd_arr); j++ {
			tpath := NewLinuxPath(false) //??
			tpath.SetUrl(dnd_arr[j])
			from_url = append(from_url, tpath)
		}
		RunFileOperaion(from_url, dest, oper)
		Prln("d&d: received [" + dest.GetReal() + "]" + dnd_str)
		//space.SetText("[" + filename + "]" + dnd_str)
	})
}
