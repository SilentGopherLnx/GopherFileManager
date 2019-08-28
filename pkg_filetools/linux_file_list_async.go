package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	//	"bufio"
	//	"io"
	"os"
	//	"os/exec"
)

const CHAN_BUFFER_SIZE = 100
const LINUX_SMB = "smb://"
const LINUX_DISKS = "pc://"

func init() {

}

type LinuxFileReport struct {
	IsVirtual   bool
	NameOnly    string
	FullName    *LinuxPath
	IsRegular   bool
	IsDirectory bool
	IsLink      bool
	SizeBytes   int64
	Mode        string
	ModTime     Time
}

type IFileListAsync interface {
	CheckUpdatesFinished() ([]*LinuxFileReport, bool, error)
	Kill()
}

type FileListAsync struct {
	path            string
	method          IFileListAsync
	time_last       Time
	period_constant float64
	finished        bool
	all_data        []*LinuxFileReport
	err             error
}

func NewFileListAsync(path_real string, notify_period float64, method IFileListAsync) *FileListAsync {
	fla := FileListAsync{
		path:            path_real,
		method:          method,
		time_last:       TimeNow(),
		period_constant: notify_period,
	}
	return &fla
}

//have news, data, finish, errors
func (w *FileListAsync) CheckUpdatesFinished() (bool, []*LinuxFileReport, bool, error) {
	RuntimeGosched()
	if !w.finished {
		time_now := TimeNow()
		if TimeSecondsSub(w.time_last, time_now) > w.period_constant {
			w.time_last = time_now
			data, finish, err := w.method.CheckUpdatesFinished()
			if finish {
				w.finished = true
			}
			if err != nil {
				w.err = err
				w.finished = true
			}
			if len(data) > 0 {
				w.all_data = append(w.all_data, data...)
				return true, data, finish, w.err
			} else {
				return false, []*LinuxFileReport{}, finish, w.err
			}
		}
	}
	return false, []*LinuxFileReport{}, w.finished, w.err
}

func (w *FileListAsync) ForceKill() {
	w.method.Kill()
	w.finished = true
}

func (w *FileListAsync) AllData() ([]*LinuxFileReport, error) {
	return w.all_data, w.err
}

//========================
//========================

func NewFileListAsync_DetectType(path *LinuxPath, search_name string, buffer_size int, notify_period float64) *FileListAsync {
	buffer := MAXI(1, buffer_size)
	url := path.GetUrl()
	len_smb := StringLength(LINUX_SMB)
	if url == LINUX_SMB {
		// if search_name != "" {
		// 	return nil
		// }
		return NewFileListAsync_NetworkAll(notify_period)
	}
	if url == LINUX_DISKS {
		// if search_name != "" {
		// 	return nil
		// }
		return NewFileListAsync_Disks(notify_period)
	}
	if StringPart(url, 1, len_smb) == LINUX_SMB {
		pc_name := StringPart(url, len_smb+1, 0)
		pc_name = FilePathEndSlashRemove(pc_name)
		if StringFind(pc_name, GetOS_Slash()) == 0 {
			// if search_name != "" {
			// 	return nil
			// }
			return NewFileListAsync_NetworkFolders(pc_name, notify_period)
		}
	}
	if search_name != "" {
		return NewFileListAsync_Searcher(path.GetReal(), search_name, notify_period)
	} else {
		// if StringFind(url, "mtp://") == 1 || StringFind(url, "gphoto2://") == 1 ||
		// 	StringFind(url, "dav://") == 1 || StringFind(url, "davs://") == 1 ||
		// 	StringFind(url, "ftp://") == 1 || StringFind(url, "ftps://") == 1 {
		// 	return NewFileListAsync_Directory(path.GetReal(), buffer, notify_period)
		// }
		// return NewFileListAsync_Directory(path.GetReal(), -1, notify_period)
		if StringFind(url, "file://") == 1 || StringFind(url, "smb://") == 1 {
			return NewFileListAsync_Directory(path.GetReal(), -1, notify_period)
		} else {
			return NewFileListAsync_Directory(path.GetReal(), buffer, notify_period)
		}
	}
}

type fileListAsync_ struct {
	lock *SyncMutex
	data []*LinuxFileReport
	done int
	err  error
}

func (m *fileListAsync_) checkUpdatesFinished_() ([]*LinuxFileReport, bool, error) {
	m.lock.Lock()
	// defer m.lock.Unlock()
	if m.done < 2 {
		if len(m.data) == 0 && m.done == 1 {
			m.done = 2
			m.lock.Unlock()
			return []*LinuxFileReport{}, true, nil
		} else {
			data := m.data
			m.data = []*LinuxFileReport{}
			m.lock.Unlock()
			return data, false, nil
		}
	} else {
		m.lock.Unlock()
		return []*LinuxFileReport{}, true, nil
	}
}

func (m *fileListAsync_) kill_() {
	m.lock.Lock()
	m.done = 2
	m.lock.Unlock()
}

//========================
//========================

func NewFileListAsync_Directory(path_real string, buffer_size int, notify_period float64) *FileListAsync {
	m := &fileListAsync_dir{}
	m.lock = NewSyncMutex()
	m.kill = NewAtomicBool(false)
	// m.cmd = exec.Command("ls", "-U", "-a", "-p", "-1")
	// m.cmd.Dir = path_real
	chan_found := make(chan *LinuxFileReport, CHAN_BUFFER_SIZE)
	bs := buffer_size
	if buffer_size < 1 {
		bs = -1
	}
	go func() {
		f, err0 := os.Open(path_real)
		if err0 == nil {
			//slash := GetOS_Slash()
			for {
				if m.kill.Get() {
					close(chan_found)
					break
				} else {
					files, err1 := f.Readdir(bs)
					if err1 == nil {
						for j := 0; j < len(files); j++ {
							//linestr := files[j].Name() + B2S(files[j].IsDir(), slash, "")

							f := &LinuxFileReport{IsVirtual: false}
							f.NameOnly = files[j].Name()
							f.FullName = NewLinuxPath(files[j].IsDir())
							f.FullName.SetReal(FolderPathEndSlash(path_real) + files[j].Name())
							f.IsRegular = files[j].Mode().IsRegular()
							f.IsDirectory = files[j].IsDir()
							f.IsLink = FileIsLink(files[j])
							f.SizeBytes = files[j].Size()
							f.Mode = files[j].Mode().String()
							f.ModTime = Time(files[j].ModTime())

							chan_found <- f
						}
						if bs == -1 {
							close(chan_found)
							break
						}
					} else {
						close(chan_found)
						break
					}
				}
			}
			f.Close()
		}
	}()
	go func() {
		for {
			name, ok := <-chan_found
			m.lock.Lock()
			if ok {
				m.data = append(m.data, name)
				m.lock.Unlock()
			} else {
				m.done = 1
				m.lock.Unlock()
				break
			}
			//SleepMS(5)
			//RuntimeGosched()
		}
	}()
	return NewFileListAsync(path_real, notify_period, m)
}

type fileListAsync_dir struct {
	fileListAsync_
	kill *ABool
}

func (m *fileListAsync_dir) CheckUpdatesFinished() ([]*LinuxFileReport, bool, error) {
	return m.checkUpdatesFinished_()
}

func (m *fileListAsync_dir) Kill() {
	// if m.cmd != nil && m.cmd.Process != nil {
	// 	m.lock.Lock()
	// 	m.cmd.Process.Kill()
	// 	m.lock.Unlock()
	// }
	m.kill.Set(true)
	m.kill_()
}

//========================

func NewFileListAsync_Searcher(path_real string, search_name string, notify_period float64) *FileListAsync {
	mount_list := LinuxGetMountList()
	dir, err := FileInfo(path_real, true)
	if err == nil && dir.IsDir() {
		kill := NewAtomicBool(false)
		m := &fileListAsync_search{kill: kill}
		m.lock = NewSyncMutex()
		chan_found := make(chan *FileReport, CHAN_BUFFER_SIZE)
		go FoldersRecursively_Search(mount_list, dir, path_real, search_name, chan_found, kill)
		go func() {
			for {
				r, ok := <-chan_found
				m.lock.Lock()
				if ok {
					f := &LinuxFileReport{IsVirtual: false}
					f.NameOnly = r.NameOnly
					f.FullName = NewLinuxPath(r.IsDir())
					f.FullName.SetReal(r.FullName)
					f.IsRegular = r.IsRegular
					f.IsDirectory = r.IsDir()
					f.IsLink = r.IsLink
					f.SizeBytes = r.Size()
					f.Mode = r.Mode().String()
					f.ModTime = Time(r.ModTime())

					m.data = append(m.data, f)
					m.lock.Unlock()
				} else {
					m.done = 1
					m.lock.Unlock()
					break
				}
				//SleepMS(5)
				RuntimeGosched()
			}
		}()
		return NewFileListAsync(path_real, notify_period, m)
	}
	return nil
}

type fileListAsync_search struct {
	fileListAsync_
	kill *ABool
}

func (m *fileListAsync_search) CheckUpdatesFinished() ([]*LinuxFileReport, bool, error) {
	m.lock.Lock()
	if m.done < 2 {
		if len(m.data) == 0 && m.done == 1 {
			m.done = 2
			m.lock.Unlock()
			return []*LinuxFileReport{}, true, nil
		} else {
			data := m.data
			m.data = []*LinuxFileReport{}
			m.lock.Unlock()
			return data, false, nil
		}
	} else {
		m.lock.Unlock()
		return []*LinuxFileReport{}, true, nil
	}
}

func (m *fileListAsync_search) Kill() {
	m.kill.Set(true)
	m.kill_()
}

//========================

func NewFileListAsync_NetworkAll(notify_period float64) *FileListAsync {
	m := &fileListAsync_netall{}
	m.lock = NewSyncMutex()
	m.data = []*LinuxFileReport{}
	go func() {
		smbs, err := SMB_ScanNetwork()
		data := []*LinuxFileReport{}
		if err == nil {
			for j := 0; j < len(smbs); j++ {
				f := &LinuxFileReport{IsVirtual: true}
				f.NameOnly = FilePathEndSlashRemove(StringDown(smbs[j].Name))
				f.FullName = NewLinuxPath(true)
				f.FullName.SetUrl(LINUX_SMB + f.NameOnly + "/")
				f.IsRegular = false
				f.IsDirectory = true
				f.IsLink = false
				f.SizeBytes = 0
				f.Mode = ""
				f.ModTime = TimeZero()

				data = append(data, f)
			}
		}
		m.lock.Lock()
		m.done = 1
		if err == nil {
			m.data = data
		} else {
			m.err = err
		}
		m.lock.Unlock()

	}()
	return NewFileListAsync(LINUX_SMB, notify_period, m)
}

type fileListAsync_netall struct {
	fileListAsync_
}

func (m *fileListAsync_netall) CheckUpdatesFinished() ([]*LinuxFileReport, bool, error) {
	return m.checkUpdatesFinished_()
}

func (m *fileListAsync_netall) Kill() {
	m.kill_()
}

//========================

func NewFileListAsync_NetworkFolders(pc_name string, notify_period float64) *FileListAsync {
	m := &fileListAsync_netall{}
	m.lock = NewSyncMutex()
	m.data = []*LinuxFileReport{}
	go func() {
		folders, err := SMB_GetPublicFolders(pc_name)
		data := []*LinuxFileReport{}
		if err == nil {
			for j := 0; j < len(folders); j++ {
				f := &LinuxFileReport{IsVirtual: true}
				f.NameOnly = FilePathEndSlashRemove(folders[j])
				f.FullName = NewLinuxPath(true)
				f.FullName.SetUrl(LINUX_SMB + pc_name + "/" + f.NameOnly + "/")
				f.IsRegular = false
				f.IsDirectory = true
				f.IsLink = false
				f.SizeBytes = 0
				f.Mode = ""
				f.ModTime = TimeZero()

				data = append(data, f)
			}
		}
		m.lock.Lock()
		m.done = 1
		if err == nil {
			m.data = data
		} else {
			m.err = err
		}
		m.lock.Unlock()

	}()
	return NewFileListAsync(LINUX_SMB+pc_name+"/", notify_period, m)
}

type fileListAsync_netfolder struct {
	fileListAsync_
}

func (m *fileListAsync_netfolder) CheckUpdatesFinished() ([]*LinuxFileReport, bool, error) {
	return m.checkUpdatesFinished_()
}

func (m *fileListAsync_netfolder) Kill() {
	m.kill_()
}

//========================

func NewFileListAsync_Disks(notify_period float64) *FileListAsync {
	m := &fileListAsync_netall{}
	m.lock = NewSyncMutex()
	m.data = []*LinuxFileReport{}
	go func() {
		mounts := Linux_DisksGetAllLocal()
		data := []*LinuxFileReport{}
		for j := 0; j < len(mounts); j++ {
			f := &LinuxFileReport{IsVirtual: true}
			f.NameOnly = FilePathEndSlashRemove(mounts[j].Title)
			f.FullName = NewLinuxPath(true)
			f.FullName.SetReal(mounts[j].MountPath)
			f.IsRegular = false
			f.IsDirectory = true
			f.IsLink = false
			f.SizeBytes = 0
			f.Mode = ""
			f.ModTime = TimeZero()

			data = append(data, f)
		}
		m.lock.Lock()
		m.done = 1
		m.data = data
		m.lock.Unlock()

	}()
	return NewFileListAsync(LINUX_DISKS, notify_period, m)
}

type fileListAsync_disks struct {
	fileListAsync_
}

func (m *fileListAsync_disks) CheckUpdatesFinished() ([]*LinuxFileReport, bool, error) {
	return m.checkUpdatesFinished_()
}

func (m *fileListAsync_disks) Kill() {
	m.kill_()
}
