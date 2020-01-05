package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"

	"github.com/fsnotify/fsnotify" // golang.org/x/sys/unix
)

type FSWatcher struct {
	fswatcher       *fsnotify.Watcher
	path            string
	period_constant float64
	count_main      *AInt
	time_last_main  Time
	//count_write     *AInt
	time_last_write Time
	lock_write      *SyncMutex
	was_write       bool
}

func NewFSWatcher(notify_period int) *FSWatcher {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		Prln(err.Error())
		fswatcher = nil
	}
	w := FSWatcher{fswatcher: fswatcher, count_main: NewAtomicInt(0), time_last_main: TimeNow(), time_last_write: TimeNow(), lock_write: NewSyncMutex()} //, count_write: NewAtomicInt(0)
	w.period_constant = float64(notify_period)
	return &w
}

func (w *FSWatcher) SetListenerOnce() {
	//done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-w.fswatcher.Events:
				if ok {
					evname := StringUp(StringEnd(StringTrim(event.String()), 5))
					if evname != "WRITE" {
						//if event.Op&fsnotify.Write != fsnotify.Write { //not work
						Prln("FS NOTIFY EVENT:" + event.String()) // + "[" + evname + "][" + event.Name + "]")
						w.count_main.Add(1)
					} else {
						//Prln("FS NOTIFY - modified file: " + event.Name)
						//Prln("FS NOTIFY EVENT: [WRITE]")
						w.lock_write.Lock()
						w.time_last_write = TimeNow()
						//w.count_write.Add(1)
						w.was_write = true
						w.lock_write.Unlock()
					}
				}
			case err, ok := <-w.fswatcher.Errors:
				if ok {
					Prln("FS NOTIFY ERROR: " + err.Error())
				}
			}
		}
	}()
}

func (w *FSWatcher) IsUpdated() bool {
	v := false
	time_now := TimeNow()
	if TimeSecondsSub(w.time_last_main, time_now) > w.period_constant {
		n := w.count_main.Get()
		if n > 0 {
			w.time_last_main = time_now
			w.count_main.Set(0)
			v = true
		}
		w.lock_write.Lock()
		if TimeSecondsSub(w.time_last_write, time_now) > w.period_constant*1.5 {
			//if w.count_write.Get() > 0 {
			if w.was_write {
				w.was_write = false
				//w.count_write.Set(0)
				v = true
			}
		}
		w.lock_write.Unlock()
	}
	return v
}

func (w *FSWatcher) Select(path string) {
	path2 := FilePathEndSlashRemove(path)
	if w.path != path2 {
		if w.fswatcher != nil {
			w.fswatcher.Remove(w.path)
			//w.count.Set(0)
			w.path = path2
			w.fswatcher.Add(w.path)
		} else {
			w.path = path
		}
	}
}

func (w *FSWatcher) Close() {
	w.Close()
}

func (w *FSWatcher) EmulateUpdate() {
	w.count_main.Add(1)
}
