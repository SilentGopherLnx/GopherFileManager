package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	//	. "github.com/SilentGopherLnx/GopherFileManager/pkg_fileicon"

	"github.com/gotk3/gotk3/gtk"
	//	"github.com/gotk3/gotk3/gdk"
)

var rwlock *SyncMutexRW_OneWriterProtected = NewSyncMutexRW_OneWriterProtected()

func Thread_Main() {
	iter := 0
	gtk.MainIteration()
	RuntimeGosched()
	for {
		if async != nil {
			updated, data, finish, err := async.CheckUpdatesFinished()
			if updated && !finish {
				lpath2 := FolderPathEndSlash(path.GetReal())
				AddFilesToList(gGFiles, data, lpath2, req_id.Get())
			}
			if finish {
				err = err
				spinnerFiles.Stop()
				if err != nil {
					Dialog_FolderError(win, err, path.GetVisual())
				}
				async = nil
				rwlock.W_Unlock()
			}
		}

		if fswatcher.IsUpdated() {
			listFiles(gGFiles, path, false, false)
		}
		gtk.MainIteration()
		qlen := qu.Length()
		if qlen > 0 {
			//Prln("qlen:" + I2S(qlen) + " / " + F2S(GetPC_MemoryUsageMb(), 1) + "Mb")
			//Prln("it1")
			w, ok := qu.GetEnd().(*IconUpdateable)
			for ok && !GTK_WidgetExist(w.widget.GetIcon()) && qu.Length() > 0 {
				w, ok = qu.GetEnd().(*IconUpdateable)
				Prln("widget searching...")
			}
			if ok && GTK_WidgetExist(w.widget.GetIcon()) {
				//Prln("pixbufset")
				if w.success_preview {
					w.widget.SetLoading(false, false)
					w.widget.GetIcon().SetFromPixbuf(w.pixbuf_preview)
				} else {
					if w.pixbuf_cache != nil {
						w.widget.GetIcon().SetFromPixbuf(w.pixbuf_cache)
					} else {
						w.widget.SetLoading(false, true)
					}
				}
			}
			//Prln("it2")
		} else {
			iter++
		}
		//if iter > 10 {
		//	iter = 0
		mem.SetText(I2S(num_works.Get()) + " " + langs.GetStr("gui_down_processes") + "; RAM Usage: " + F2S(GetPC_MemoryUsageMb(), 1) + " Mb & " + usage + "; displayed:" + I2S(len(arr_blocks)) + "/" + I2S(real_files_count) + " objects")
		main_iterations_funcs.ExecAll()

		if num_works.Get() == 0 {
			if GTK_SpinnerActive(spinnerIcons, true) {
				spinnerIcons.Stop()
			}
		} else {
			if !GTK_SpinnerActive(spinnerIcons, false) {
				spinnerIcons.Start()
			}
		}

		cb, cf := hist.CanBackForward()
		gBtnBack.SetSensitive(cb)
		gBtnForward.SetSensitive(cf)

		//}
		//RuntimeGosched()
		//debug.FreeOSMemory()
		//mem.SetText("RAM Usage: " + I2S(linux.LinuxMemory()) + " Mb")
		//GarbageCollection()
		//win.ShowAll()
	}
}

func Thread_Icon(icon_chan chan *IconUpdateable, qu *SyncQueue, thread_id int) { // qu *queue.Queue
	Prln("thread #" + I2S(thread_id) + " has go-id: " + I2S(GoId()))
	for {
		msg := <-icon_chan
		num_works.Add(1)
		if req_id.Get() == msg.req && GTK_WidgetExist(msg.widget.GetIcon()) {
			//var pixbuf_preview *gdk.Pixbuf
			var ok = false
			var skip = false
			if msg.skip_if_cache_loaded && msg.pixbuf_cache != nil {
				skip = true
			}
			if msg.folder {
				//if !msg.basic_mode {
				msg.pixbuf_preview, ok = GetPixBufGTK_Folder(msg.fullname, ZOOM_SIZE, msg.basic_mode, qu, nil, msg.req, skip) //, msg)
				//}
			} else {
				if skip {
					msg.pixbuf_preview = msg.pixbuf_cache
					ok = true
				} else {
					if StringInArray(msg.tfile, MIME_IMAGE) > -1 {
						msg.pixbuf_preview, ok = GetPreview_ImagePixBuf(msg.fullname, ZOOM_SIZE)
					}
					if !msg.basic_mode && StringInArray(msg.tfile, MIME_VIDEO) > -1 {
						msg.pixbuf_preview, ok = GetPreview_VideoPixBuf(msg.fullname, ZOOM_SIZE, nil, msg.req, true)
					}
				}
			}

			if ok {
				msg.success_preview = true
				qu.Append(msg)
			} else {
				qu.Append(msg)
			}
		}
		num_works.Add(-1)
		RuntimeGosched()
	}
}

func Thread_GC_and_Free(pid int) {
	for {
		SleepMS(1500)
		GarbageCollection()
		FreeOSMemory()
		usage = F2S(LinixMemoryUsedMB(pid), 1) + "Mb"
	}
}
