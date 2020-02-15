package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var select_mode bool
var select_x1, select_y1, select_x2, select_y2 int

var select_with_old bool
var selected_old []int

func FilesSelector_GetList() []string {
	arr := []string{}
	for j := 0; j < len(arr_blocks); j++ {
		if arr_blocks[j].GetSelected() {
			name := arr_blocks[j].GetFileName()
			if arr_blocks[j].IsDir() {
				arr = append(arr, FolderPathEndSlash(name))
			} else {
				arr = append(arr, name)
			}
		}
	}
	return arr
}

func FilesSelector_Draw(dy int, ctx *cairo.Context) {
	if select_x1 > 0 && select_y1 > 0 && select_x2 > 0 && select_y2 > 0 {
		c := GTK_ColorOfSelected()
		ctx.SetSourceRGBA(c[0], c[1], c[2], 1.0) //0.4, 0.7, 0.8, 1.0) // BLUE DARK
		ctx.Rectangle(float64(select_x1), float64(select_y1-dy), float64(select_x2-select_x1), float64(select_y2-select_y1))
		ctx.Fill()
	}
}

func FilesSelector_MouseAtSelectZone(x0, y0 int) bool {
	at_zone := false
	for j := 0; j < len(arr_blocks); j++ {
		at_zone = at_zone || arr_blocks[j].IsClickedIn(&gGFiles.Widget, x0, y0)
	}
	return !at_zone
}

func FilesSelector_MousePressed(event *gdk.Event, scroll *gtk.ScrolledWindow) (int, int, int, bool) {
	mousekey, x1, y1, _ := GTK_MouseKeyOfEvent(event)
	_, dy := GTK_ScrollGetValues(scroll)
	y1 += dy
	zone := FilesSelector_MouseAtSelectZone(x1, y1)
	if mousekey == 1 && zone {
		select_x1 = x1
		select_y1 = y1
		select_x2 = 0
		select_y2 = 0
		Prln("select mouse1_down " + I2S(x1) + "/" + I2S(y1+dy))
		//scroll.GrabFocus()
	}
	selected_old = []int{}
	for j := 0; j < len(arr_blocks); j++ {
		if arr_blocks[j].GetSelected() {
			selected_old = append(selected_old, j)
		}
	}
	return mousekey, x1, y1, zone
}

func FilesSelector_MouseMoved(event *gdk.Event, scroll *gtk.ScrolledWindow, redraw func()) {
	if select_x1 > 0 && select_y1 > 0 {
		_, x2, y2, _ := GTK_MouseKeyOfEvent(event)
		_, dy := GTK_ScrollGetValues(scroll)
		y2 += dy
		select_x2 = x2
		select_y2 = y2
		//Prln("rect " + I2S(select_x1) + "," + I2S(select_y1) + " / " + I2S(select_x2) + "," + I2S(select_y2))
		for j := 0; j < len(arr_blocks); j++ {
			is_inside := arr_blocks[j].IsInSelectRect(&gGFiles.Widget, select_x1, select_y1, select_x2, select_y2)
			arr_blocks[j].SetSelected(is_inside || IntInArray(j, selected_old) > -1)
		}
		redraw()
	}
}

func FilesSelector_MouseRelease(event *gdk.Event, scroll *gtk.ScrolledWindow, redraw func()) {
	mousekey, x2, y2, _ := GTK_MouseKeyOfEvent(event)
	_, dy := GTK_ScrollGetValues(scroll)
	y2 += dy
	zone2 := FilesSelector_MouseAtSelectZone(x2, y2)
	if mousekey == 1 {
		if select_x1 == x2 && select_y1 == y2 && zone2 {
			Prln("select mouse1_up with reset")
			FilesSelector_ResetChecks()
		} else {
			Prln("select mouse1_up")
		}
		FilesSelector_ResetRect()
		redraw()
		selected_old = []int{}
	}
}

func FilesSelector_ResetRect() {
	select_x1 = 0
	select_y1 = 0
	select_x2 = 0
	select_y2 = 0
}

func FilesSelector_ResetChecks() {
	for j := 0; j < len(arr_blocks); j++ {
		arr_blocks[j].SetSelected(false)
	}
}

func FilesSelector_SelectAll() {
	for j := 0; j < len(arr_blocks); j++ {
		arr_blocks[j].SetSelected(true)
	}
}
