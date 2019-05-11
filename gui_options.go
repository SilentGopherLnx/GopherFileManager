package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	"github.com/gotk3/gotk3/gtk"
)

type OptionsContainer struct {
	st *OptionsStorage
}

var opt *OptionsContainer = nil

const OPTIONS_FILE = "options.cfg"

const OPTIONS_INIT_ZOOM = "init_zoom"
const OPTIONS_APP_FILEMOVER = "filemover_path"
const OPTIONS_MOVER_BUFFER = "filemover_buffer"
const OPTIONS_FOLDER_HASH = "hash_folderpath"
const OPTIONS_NUM_THREADS = "hash_num_threads"
const OPTIONS_SYSTEM_TERMINAL = "system_terminal"
const OPTIONS_SYSTEM_FILEMANAGER = "system_filemanager"
const OPTIONS_SYSTEM_TEXTEDITOR = "system_texteditor"
const OPTIONS_FFMPEG_TIMEOUT = "ffmpeg_timeout"
const OPTIONS_INOTIFY_PERIOD = "inotify_period"
const OPTIONS_SYMLINKS_EVAL = "symlinks_eval"

func init() {
	InitOptions()
}

func InitOptions() {
	if opt == nil {
		opt = &OptionsContainer{st: NewOptionsStorage()}
	} else {
		return
	}

	opt.st.AddRecord_Array(100, OPTIONS_INIT_ZOOM, "128", []string{"64", "128", "256"}, "Initial zoom")

	opt.st.AddRecord_String(201, OPTIONS_APP_FILEMOVER, "FileMoverGui", "FileMover's App path (relative or absolute)")
	opt.st.AddRecord_Array(202, OPTIONS_MOVER_BUFFER, "16", []string{"1", "4", "8", "16", "32", "64", "128"}, "File mover buffer size in bytes (multiplied by 1024)")

	opt.st.AddRecord_String(301, OPTIONS_FOLDER_HASH, "hash/", "Path (relative or absolute) to hash folder (should exist)")
	opt.st.AddRecord_Array(302, OPTIONS_NUM_THREADS, "12", []string{"1", "2", "4", "6", "8", "12", "16", "24", "32"}, "Number of thrumbnails icon threads (need restart)")

	opt.st.AddRecord_String(401, OPTIONS_SYSTEM_TERMINAL, "gnome-terminal --working-directory=%F", "System terminal app for run. %F - path argument")
	opt.st.AddRecord_String(402, OPTIONS_SYSTEM_FILEMANAGER, "nemo %F", "System file manager app for run. %F - path argument")
	opt.st.AddRecord_String(403, OPTIONS_SYSTEM_TEXTEDITOR, "xed", "Your text editor (without arguments)")

	opt.st.AddRecord_Integer(501, OPTIONS_FFMPEG_TIMEOUT, 6, 2, 10, "ffmpeg max time wait before kill it")
	opt.st.AddRecord_Integer(502, OPTIONS_INOTIFY_PERIOD, 2, 1, 5, "inotify minimum period of directory content change reaction")

	opt.st.AddRecord_Boolean(601, OPTIONS_SYMLINKS_EVAL, true, "Open symlinks as real path to folder")

	opt.st.RecordsValues_Load(FolderLocation_App() + OPTIONS_FILE)
	opt.st.RecordsValues_Save(FolderLocation_App() + OPTIONS_FILE)

}

func (o *OptionsContainer) GetZoom() int {
	z := o.st.ValueGetString(OPTIONS_INIT_ZOOM)
	return S2I(z)
}

func (o *OptionsContainer) GetThreads() int {
	z := o.st.ValueGetString(OPTIONS_NUM_THREADS)
	return S2I(z)
}

func (o *OptionsContainer) GetMoverBuffer() int {
	z := o.st.ValueGetString(OPTIONS_MOVER_BUFFER)
	return S2I(z)
}

func (o *OptionsContainer) GetFileMover() string {
	app := o.st.ValueGetString(OPTIONS_APP_FILEMOVER)
	if StringPart(app, 1, 1) != "/" {
		app = FolderLocation_App() + app
	}
	return app
}

func (o *OptionsContainer) GetHashFolder() string {
	folder := o.st.ValueGetString(OPTIONS_FOLDER_HASH)
	if StringPart(folder, 1, 1) != "/" {
		folder = FolderLocation_App() + folder
	}
	return FolderPathEndSlash(folder)
}

func (o *OptionsContainer) GetTerminal(path string) string {
	return StringReplace(o.st.ValueGetString(OPTIONS_SYSTEM_TERMINAL), "%F", ExecQuote(path))
}

func (o *OptionsContainer) GetFileManager(path string) string {
	return StringReplace(o.st.ValueGetString(OPTIONS_SYSTEM_FILEMANAGER), "%F", ExecQuote(path))
}

func (o *OptionsContainer) GetTextEditor() string {
	return o.st.ValueGetString(OPTIONS_SYSTEM_TEXTEDITOR)
}

func (o *OptionsContainer) GetInotifyPeriod() int {
	return o.st.ValueGetInteger(OPTIONS_INOTIFY_PERIOD)
}

func (o *OptionsContainer) GetFfmpegTimeout() int {
	return o.st.ValueGetInteger(OPTIONS_FFMPEG_TIMEOUT)
}

func (o *OptionsContainer) GetSymlinkEval() bool {
	return o.st.ValueGetBoolean(OPTIONS_SYMLINKS_EVAL)
}

func Dialog_Options(w *gtk.Window) {
	winw, winh := 650, 400
	win2, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	win2.SetTitle("Options")
	win2.SetDefaultSize(winw, winh)
	win2.SetPosition(gtk.WIN_POS_CENTER)
	win2.SetTransientFor(w)
	//win2.SetIconFromFile(FolderLocation_App() + "gui/icon.png")
	win2.SetModal(true)
	win2.SetKeepAbove(true)

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	arr := opt.st.GetRecordsKeys()
	for j := 0; j < len(arr); j++ {
		key := arr[j]
		opt_widget := GTK_OptionsWidget(opt.st, key, nil)
		if opt_widget != nil {
			opt_title, _ := gtk.LabelNew(opt.st.GetRecordComment(key))
			opt_title.SetHAlign(gtk.ALIGN_END)
			grid.Attach(opt_title, 0, j, 1, 1)
			grid.Attach(opt_widget, 1, j, 1, 1)
		}
	}

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetVExpand(true)
	scroll.SetHExpand(true)
	scroll.Add(grid)

	win2.Add(scroll)

	win2.Connect("destroy", func() {
		opt.st.RecordsValues_Save(FolderLocation_App() + OPTIONS_FILE)
	})

	win2.ShowAll()
	win2.SetSizeRequest(winw, winh)
}
