package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_filetools"

	//"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func listDiscs(g *gtk.Box) {

	GTK_Childs(g, true, true) //arrd :=
	//Prln("disc_child_len:" + I2S(len(arrd)))

	discs := Linux_DisksGetWithBookmarks(true, true, true, true)
	mountlist = LinuxGetMountList()

	for j, di := range discs {
		d := di
		Prln(d.String())
		p := d.MountPath
		gDBtn, _ := gtk.ButtonNew() //WithLabel(d.Title + "\n" + d.Model + "\n" + d.PartName)
		gDBtn.SetHExpand(true)
		//gDBtn.SetSizeRequest(LEFT_PANEL_SIZE, 32)
		//gDBtn.SetHAlign(gtk.ALIGN_CENTER)
		//gDBtn.SetJustify(gtk.JUSTIFY_CENTER)
		gDBtn.Connect("clicked", func() {
			//tpath := NewLinuxPath(true)
			path.SetReal(p)
			if StringPart(p, 1, 1) == "_" {
				path.SetVisual(StringPart(p, 2, 0))
			}
			gInpPath.SetText(path.GetVisual())
			gInpSearch.SetText("")
			listFiles(gGFiles, path, true, true)
		})
		gDBtn.SetHExpand(false)

		grid, _ := gtk.GridNew()
		grid.SetHExpand(true)

		lbl1, _ := gtk.LabelNew(d.Title)
		lbl1.SetJustify(gtk.JUSTIFY_LEFT)
		lbl1.SetHAlign(gtk.ALIGN_START)
		GTK_LabelWrapMode(lbl1, 1)
		lbl1.SetLineWrap(true)
		grid.Attach(lbl1, 0, 0, 1, 1)
		lbl1.SetMarkup("<b><u>" + HtmlEscape(d.Title) + "</u></b>")

		lbl2, _ := gtk.LabelNew(d.SpaceTotal)
		lbl2.SetJustify(gtk.JUSTIFY_RIGHT)
		lbl2.SetHAlign(gtk.ALIGN_END)
		//lbl2.SetHExpand(true)
		grid.Attach(lbl2, 1, 0, 1, 1)

		if j == 0 {
			lbl1.SetMarkup("<b><u>" + HtmlEscape(StringUp(d.PartName)) + "</u></b>")
			lbl2.SetText("HOME")
		}

		extra := StringFind(StringDown(d.PartName), "gvfsd-fuse") > 0
		if len(StringTrim(d.Model)) > 0 && !extra {
			lbl3, _ := gtk.LabelNew(d.Model)
			lbl3.SetJustify(gtk.JUSTIFY_LEFT)
			lbl3.SetHAlign(gtk.ALIGN_START)
			lbl3.SetHExpand(true)
			GTK_LabelWrapMode(lbl3, 1)
			grid.Attach(lbl3, 0, 1, 2, 1)
		}

		if j > 1 && d.SpacePercent > -1 {
			lbl4, _ := gtk.LabelNew(d.PartName)
			if extra {
				lbl4.SetText(d.Model)
			}
			lbl4.SetJustify(gtk.JUSTIFY_LEFT)
			lbl4.SetHAlign(gtk.ALIGN_START)
			lbl4.SetHExpand(true)
			GTK_LabelWrapMode(lbl4, 1)
			grid.Attach(lbl4, 0, 2, 1, 1)

			fs := d.Protocol
			if d.Protocol == "PART" {
				fs = StringUp(d.FSType)
			}
			lbl5, _ := gtk.LabelNew(fs)
			lbl5.SetMarkup("<u>" + HtmlEscape(fs) + "</u>")
			lbl5.SetJustify(gtk.JUSTIFY_RIGHT)
			lbl5.SetHAlign(gtk.ALIGN_END)
			//lbl4.SetHExpand(true)
			grid.Attach(lbl5, 1, 2, 1, 1)

			progr, _ := gtk.LevelBarNew() // ProgressBarNew()
			//progr.SetFraction(float64(d.SpacePercent) / 100.0)
			progr.SetValue(float64(d.SpacePercent) / 100.0)
			progr.SetHExpand(true)
			progr.SetSizeRequest(10, 0)

			grid.Attach(progr, 0, 3, 2, 1)
		}

		gDBtn.Connect("button-press-event", func(_ *gtk.Button, event *gdk.Event) {
			mousekey, _, _, _ := GTK_MouseKeyOfEvent(event)
			switch mousekey {
			case 3:
				Prln("right")
				rightmenu, _ := gtk.MenuNew()

				if !d.Primary && d.SpacePercent >= 0 {
					GTK_MenuItem(rightmenu, "Eject", nil)
				}
				if d.SpacePercent < 0 {
					GTK_MenuItem(rightmenu, "Remove", nil)
				}
				GTK_MenuItem(rightmenu, "Info", nil)

				rightmenu.ShowAll()
				rightmenu.PopupAtPointer(event) // (evBox, gdk.GDK_GRAVITY_STATIC, gdk.GDK_GRAVITY_STATIC,
			}
		})

		gDBtn.Add(grid)
		//g.Attach(gDBtn, 0, j, 1, 1)
		g.Add(gDBtn)
	}

	g.ShowAll()
}
