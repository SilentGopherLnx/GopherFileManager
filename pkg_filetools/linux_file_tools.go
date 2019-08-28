package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"
	//"io/ioutil"
	//"net/url"
)

var APP_EXEC_TYPE = "application/x-executable"

var SYSTEM_PATH []string

func init() {
	strs, ok := FileTextRead("system_icon_path.cfg")
	if ok {
		SYSTEM_PATH = StringSplitLines(strs)
	}
	if len(SYSTEM_PATH) < 2 {
		SYSTEM_PATH = []string{"/usr/share/icons/Mint-Y/mimetypes/64/", "/usr/share/icons/Mint-Y/apps/64/"}
	}
}

func OpenFileByApp(filename string, appname string) {
	if appname == "" {
		if FileIsExec(filename) {
			p0 := StringSplit(filename, "/")
			fname := p0[len(p0)-1]
			p1 := p0[:len(p0)-1]
			fpath := StringJoin(p1, "/")
			if len(fpath) == 0 {
				fpath = "/"
			}
			script := "cd " + ExecQuote(fpath) + " && ./" + ExecQuote(fname) + ""
			Prln("starting[" + script + "]")
			go ExecCommandBash(script)
		} else {
			safename := ExecQuote(filename)
			if FileExtension(filename) == "desktop" {
				go ExecCommandBash("`grep '^Exec' " + safename + " | tail -1 | sed 's/^Exec=//' | sed 's/%.//' | sed 's/^\"//g' | sed 's/\" *$//g'` &")
			} else {
				go ExecCommandBash("xdg-open " + safename + " &")
			}
		}
	} else {
		filetext, ok := FileTextRead("/usr/share/applications/" + appname + ".desktop")
		if ok {
			app := ""
			strs := StringSplit(filetext, "\n")
			for j := 0; j < len(strs); j++ {
				str_j := StringDown(strs[j])
				if StringFind(str_j, "exec=") == 1 {
					//Prln(str_j)
					app = StringPart(str_j, 6, 0)
				}
			}
			if len(app) > 0 {
				// %f	a single filename.
				// %F	multiple filenames.
				// %u	a single URL.
				// %U	multiple URLs.
				// %d	a single directory. Used in conjunction with %f to locate a file.
				// %D	multiple directories. Used in conjunction with %F to locate files.
				// %n	a single filename without a path.
				// %N	multiple filenames without paths.
				// %k	a URI or local filename of the location of the desktop file.
				// %v	the name of the Device entry.
				app = StringReplace(app, " %f", "")
				app = StringReplace(app, " %u", "")
				Prln("starting [" + app + "]")
				go ExecCommandBash(ExecQuote(app) + " " + ExecQuote(filename) + " &")
			} else {
				Prln("no exec in desktop file...")
			}
		} else {
			Prln("no desktop file...")
		}
	}
}

func FileMIME(filename string) string {
	ext_mime, _, _ := ExecCommand("xdg-mime", "query", "filetype", filename)
	return StringDown(StringTrim(ext_mime))
}

func AppMIME(mime_name string) string {
	ext_app, _, _ := ExecCommand("xdg-mime", "query", "default", mime_name)
	return StringTrim(ext_app)
}

//mimeopen -a file.txt
//grep 'text/x-go' -R /usr/share/applications/*
func AllAppsMIME(mime_name string) []string {
	safe_mime := StringDown(mime_name)
	desktops, _, _ := ExecCommandBash("grep " + ExecQuote(safe_mime+"=") + " /usr/share/applications/mimeinfo.cache")
	arr1 := StringSplit(desktops, ";")
	arr2 := []string{}
	desk := ".desktop"
	for j := 0; j < len(arr1); j++ {
		str := StringTrim(arr1[j])
		if StringFind(str, safe_mime+"=") == 1 {
			str = StringPart(str, StringLength(safe_mime)+2, 0)
		}
		if StringEnd(str, StringLength(desk)) == desk {
			str = StringPart(str, 1, StringLength(str)-StringLength(desk))
		}
		if len(str) > 0 {
			arr2 = append(arr2, str)
		}
	}
	return arr2
}

func FileIsExec(filename string) bool {
	return FileMIME(filename) == APP_EXEC_TYPE
}

func FileIconBySystem(filename string) string {
	print_debug := false

	ext_mime := FileMIME(filename)

	if len(ext_mime) > 0 {
		fname := SYSTEM_PATH[0] + StringReplace(ext_mime, "/", "-") + ".png"
		if FileExists(fname) {
			if print_debug {
				Prln(filename + " ## " + ext_mime)
			}
			return fname
		} else {
			ext_app := AppMIME(ext_mime)
			if len(ext_app) > 0 {
				filetext, ok := FileTextRead("/usr/share/applications/" + ext_app)
				if ok {
					icon := ""
					strs := StringSplit(filetext, "\n")
					for j := 0; j < len(strs); j++ {
						str_j := StringDown(strs[j])
						if StringFind(str_j, "icon=") == 1 {
							//Prln(str_j)
							icon = StringPart(str_j, 6, 0)
						}
					}
					if len(icon) > 0 {
						icon = SYSTEM_PATH[1] + icon + ".png"
						if print_debug {
							Prln("[" + filename + " ## " + ext_mime + " ## " + ext_app + "## " + icon + "]")
							//Prln(icon)
						}
						return icon
					}
				}
				if print_debug {
					Prln("[" + filename + " ## " + ext_mime + " ## " + ext_app + "]")
				}
			}
		}
	}
	return ""
}

func FileOrFolder_New(path string, foldermode bool) string {
	fun := func(fullname string) bool {
		if foldermode {
			return FolderMake(fullname)
		} else {
			return FileMake(fullname)
		}
	}
	fname := "New File"
	if foldermode {
		fname = "New Folder"
	}
	path2 := FolderPathEndSlash(path)
	safe_name := FileFindFreeName(path2, fname, "", "")
	if len(safe_name) > 0 {
		done := fun(path2 + safe_name)
		if done {
			return safe_name
		} else {
			return ""
		}
	}
	return ""
}
