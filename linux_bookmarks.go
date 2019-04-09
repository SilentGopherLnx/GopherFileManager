package main

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
