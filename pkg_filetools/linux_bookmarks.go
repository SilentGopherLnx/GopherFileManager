package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"
)

func LinuxGetBookmarks() [][2]string {
	_, login, _ := GetPC_UserUidLoginName()
	filename := "/home/" + login + "/.config/gtk-3.0/bookmarks"
	list, ok := FileTextRead(filename)
	if !ok {
		return [][2]string{}
	}
	arr := StringSplitLines(list)
	res := [][2]string{}
	for j := 0; j < len(arr); j++ {
		arr_j := arr[j]
		ind := StringFind(arr_j, " ")
		if ind > 0 {
			path := StringPart(arr_j, 1, ind-1)
			tpath := NewLinuxPath(true)
			tpath.SetUrl(path)
			path = tpath.GetReal()
			title := StringPart(arr_j, ind, 0)
			if len(path) > 0 && len(title) > 0 {
				res = append(res, [2]string{path, title})
			}
		}
	}
	return res
}

func Linux_DisksGetWithBookmarks(local bool, remote bool, home bool, bookmarks bool) []*DiskPart {
	disks := Linux_DisksGetMounted(local, remote)

	_, login, _ := GetPC_UserUidLoginName()

	if home {
		h := &DiskPart{Title: "<HOME>", PartName: login, FSType: "", Protocol: "HOME", SpaceTotal: "", SpaceUsed: "", SpaceFree: "", SpacePercent: 0, Crypted: false, Primary: true, MountPath: FolderLocation_UserHome(), Model: GetPC()}
		disks = append([]*DiskPart{h}, disks...)
	}

	if bookmarks {
		bmk := LinuxGetBookmarks()
		for _, b := range bmk {
			bm := &DiskPart{Title: b[1], PartName: b[0], FSType: "", Protocol: "BKMRK", SpaceTotal: "", SpaceUsed: "", SpaceFree: "", SpacePercent: -1, Crypted: false, Primary: true, MountPath: b[0], Model: b[0]}
			disks = append(disks, bm)
		}
	}
	return disks
}

//check program "ls" is installed
// func checkLsExists() {
//     path, err := exec.LookPath("ls")
//     if err != nil {
//         fmt.Printf("didn't find 'ls' executable\n")
//     } else {
//         fmt.Printf("'ls' executable is in '%s'\n", path)
//     }
//}
