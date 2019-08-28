package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"

	"github.com/fsnotify/fsnotify" // golang.org/x/sys/unix
)

type FSWatcher struct {
	fswatcher       *fsnotify.Watcher
	path            string
	count           *AInt
	time_last       Time
	period_constant float64
}

func NewFSWatcher(notify_period int) *FSWatcher {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		Prln(err.Error())
		fswatcher = nil
	}
	w := FSWatcher{fswatcher: fswatcher, count: NewAtomicInt(0), time_last: TimeNow()}
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
					w.count.Add(1)
					Prln("FS NOTIFY EVENT:" + event.String())
					// if event.Op&fsnotify.Write == fsnotify.Write {
					// 	Prln("FS NOTIFY - modified file: " + event.Name)
					// }
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
	time_now := TimeNow()
	if TimeSecondsSub(w.time_last, time_now) > w.period_constant {
		n := w.count.Get()
		if n > 0 {
			w.time_last = time_now
			w.count.Add(-n)
			return true
		}
	}
	return false
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
