package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_fileicon"

	"github.com/gotk3/gotk3/gtk"
)

type OptionsContainer struct {
	st *OptionsStorage
}

var opt *OptionsContainer = nil

const OPTIONS_FILE = "options.cfg"

const OPTIONS_LANG = "language"
const OPTIONS_INIT_ZOOM = "init_zoom"
const OPTIONS_APP_FILEMOVER = "filemover_path"
const OPTIONS_MOVER_BUFFER = "filemover_buffer"
const OPTIONS_FOLDER_HASH = "hash_folderpath"
const OPTIONS_NUM_THREADS = "hash_num_threads"
const OPTIONS_SYSTEM_TERMINAL = "system_terminal"
const OPTIONS_SYSTEM_FILEMANAGER = "system_filemanager"
const OPTIONS_SYSTEM_TEXTEDITOR = "system_texteditor"
const OPTIONS_FFMPEG_TIMEOUT = "ffmpeg_timeout"
const OPTIONS_VIDEO_PREVIEW_PERCENT = "video_preview_percent"
const OPTIONS_INOTIFY_PERIOD = "inotify_period"
const OPTIONS_SYMLINKS_EVAL = "symlinks_eval"
const OPTIONS_EXIF_ROTATION = "exif_rotation"
const OPTIONS_FOLDER_LIMIT = "folder_limit"

func InitOptions() {
	if opt == nil {
		opt = &OptionsContainer{st: NewOptionsStorage()}
	} else {
		return
	}

	opt_lang := &OptionsContainer{st: NewOptionsStorage()}
	opt_lang.st.AddRecord_Array(1, OPTIONS_LANG, "en", langs.GetLangsCodes(), OPTIONS_LANG)
	opt_lang.st.RecordsValues_Load(FolderLocation_App() + OPTIONS_FILE)
	langs.SetLang(opt_lang.GetLanguage())

	opt.st.AddRecord_Array(99, OPTIONS_LANG, "en", langs.GetLangsCodes(), "Language (Need restart!)")

	zooms_str := []string{}
	za := Constant_ZoomArray()
	for j := 0; j < len(za); j++ {
		zooms_str = append(zooms_str, I2S(za[j]))
	}
	opt.st.AddRecord_Array(100, OPTIONS_INIT_ZOOM, I2S(za[len(za)/2]), zooms_str, langs.GetStr("options_init_zoom"))

	opt.st.AddRecord_String(201, OPTIONS_APP_FILEMOVER, "FileMoverGui", langs.GetStr("options_filemover_path"))
	opt.st.AddRecord_Array(202, OPTIONS_MOVER_BUFFER, "16", []string{"1", "4", "8", "16", "32", "64", "128"}, langs.GetStr("options_filemover_buffer"))

	opt.st.AddRecord_String(301, OPTIONS_FOLDER_HASH, "hash/", langs.GetStr("options_hash_path"))
	opt.st.AddRecord_Array(302, OPTIONS_NUM_THREADS, "12", []string{"1", "2", "4", "6", "8", "12", "16", "24", "32"}, langs.GetStr("options_threads"))

	opt.st.AddRecord_String(401, OPTIONS_SYSTEM_TERMINAL, "gnome-terminal --working-directory=%F", "System terminal app for run. %F - path argument")
	opt.st.AddRecord_String(402, OPTIONS_SYSTEM_FILEMANAGER, "nemo %F", "System file manager app for run. %F - path argument")
	opt.st.AddRecord_String(403, OPTIONS_SYSTEM_TEXTEDITOR, "xed", "Your text editor (without arguments)")

	opt.st.AddRecord_Integer(501, OPTIONS_FFMPEG_TIMEOUT, 6, 2, 10, langs.GetStr("options_ffmpeg"))
	opt.st.AddRecord_Integer(502, OPTIONS_INOTIFY_PERIOD, 2, 1, 5, langs.GetStr("options_inotify"))
	opt.st.AddRecord_Integer(503, OPTIONS_VIDEO_PREVIEW_PERCENT, 50, 1, 99, langs.GetStr("options_video_percent"))

	opt.st.AddRecord_Integer(504, OPTIONS_FOLDER_LIMIT, 10, 2, 50, langs.GetStr("options_max_result"))

	opt.st.AddRecord_Boolean(601, OPTIONS_SYMLINKS_EVAL, true, "Open symlinks as real path to folder")
	opt.st.AddRecord_Boolean(602, OPTIONS_EXIF_ROTATION, true, "Use EXIF orientation tag for JPEG images")

	opt.st.RecordsValues_Load(FolderLocation_App() + OPTIONS_FILE)
	opt.st.RecordsValues_Save(FolderLocation_App() + OPTIONS_FILE)

}

func (o *OptionsContainer) GetLanguage() string {
	z := o.st.ValueGetString(OPTIONS_LANG)
	if StringTrim(z) == "" {
		return DEFAULT_LANG
	}
	return z
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

func (o *OptionsContainer) GetVideoPercent() int {
	return o.st.ValueGetInteger(OPTIONS_VIDEO_PREVIEW_PERCENT)
}

func (o *OptionsContainer) GetSymlinkEval() bool {
	return o.st.ValueGetBoolean(OPTIONS_SYMLINKS_EVAL)
}

func (o *OptionsContainer) GetExifRot() bool {
	return o.st.ValueGetBoolean(OPTIONS_EXIF_ROTATION)
}

func (o *OptionsContainer) GetFolderLimit() int {
	return o.st.ValueGetInteger(OPTIONS_FOLDER_LIMIT) * 100
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

	win2.Connect("destroy", func() {
		//opt.st.RecordsValues_Save(FolderLocation_App() + OPTIONS_FILE)
	})

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	clear_hash, _ := gtk.ButtonNewWithLabel("")
	upd_hash_clear := func(key string) {
		if key == OPTIONS_FOLDER_HASH {
			files_list, _ := Folder_ListFiles(opt.GetHashFolder(), false)
			files_num := len(files_list)
			newname := "Clear Hash (" + I2S(files_num) + " files)"
			clear_hash.SetLabel(newname)
		}
	}
	upd_hash_clear(OPTIONS_FOLDER_HASH)

	arr := opt.st.GetRecordsKeys()
	arr_len := len(arr)
	for j := 0; j < arr_len; j++ {
		key := arr[j]
		opt_widget := GTK_OptionsWidget(opt.st, key, upd_hash_clear)
		if opt_widget != nil {
			opt_title, _ := gtk.LabelNew(opt.st.GetRecordComment(key))
			opt_title.SetHAlign(gtk.ALIGN_END)
			grid.Attach(opt_title, 0, j, 1, 1)
			grid.Attach(opt_widget, 1, j, 1, 1)
		}
	}
	clear_hash.Connect("button-release-event", func() {
		//Prln("" + opt.GetHashFolder())
		file1 := NewLinuxPath(false) //??
		file1.SetReal(opt.GetHashFolder())
		RunFileOperaion([]*LinuxPath{file1}, nil, OPER_CLEAR)
	})
	grid.Attach(clear_hash, 1, arr_len, 1, 1)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetVExpand(true)
	scroll.SetHExpand(true)
	scroll.Add(grid)

	win2.Add(scroll)

	win2.Connect("destroy", func() {
		opt.st.RecordsValues_Save(FolderLocation_App() + OPTIONS_FILE)
		langs.SetLang(opt.GetLanguage())
	})

	win2.ShowAll()
	win2.SetSizeRequest(winw, winh)
}
