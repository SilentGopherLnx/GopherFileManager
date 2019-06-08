package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"image"
	"image/color"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/pango"
)

var pixbuf_link, pixbuf_notread, pixbuf_mount, pixbuf_loading, pixbuf_loading_err *gdk.Pixbuf
var hid map[int]*gdk.Pixbuf

const BACK_GRAY_VISIBLE float64 = 0.8
const BACK_GRAY_HIDDEN float64 = 0.9

// //get the icon theme and lookup the icon we want by name, here at a size of 64px
// var info = Gtk.IconTheme.get_default ().lookup_icon ("view-refresh-symbolic", 64, 0);
// //now load the icon as a symbolic with a color set in the brackets as RGBA, here as plain red
// var pixbuf = info.load_symbolic ({1, 0, 0, 1});

// sudo apt-get install gtk-3-examples
// gtk3-icon-browser

func init() {
	appdir := FolderLocation_App()
	bb, _ := FileBytesRead(appdir + "gui/emblem_loading.png")
	pixbuf_loading = GTK_PixBuf_From_Bytes(bb, "png")

	bb, _ = FileBytesRead(appdir + "gui/emblem_loading_error.png")
	pixbuf_loading_err = GTK_PixBuf_From_Bytes(bb, "png")

	bb, _ = FileBytesRead(appdir + "gui/emblem_link.png")
	pixbuf_link = GTK_PixBuf_From_Bytes(bb, "png")

	bb, _ = FileBytesRead(appdir + "gui/emblem_unreadable.png")
	pixbuf_notread = GTK_PixBuf_From_Bytes(bb, "png")

	bb, _ = FileBytesRead(appdir + "gui/emblem_mount.png")
	pixbuf_mount = GTK_PixBuf_From_Bytes(bb, "png")

	gr := uint8(RoundF(float64(255) * BACK_GRAY_HIDDEN))
	zooms := Constant_ZoomArray()
	for j := 0; j < len(zooms); j++ {
		wh := zooms[j]
		img := image.NewNRGBA(image.Rect(0, 0, wh, wh))
		col := color.NRGBA{R: gr, G: gr, B: gr, A: 127}
		for y := 0; y < wh; y++ {
			for x := 0; x < wh; x++ {
				if (x+y)%2 == 0 {
					img.Set(x, y, col)
				}
			}
		}
		hid = make(map[int]*gdk.Pixbuf)
		hid[wh] = GTK_PixBuf_From_RGBA(img)
	}

	//bb, _ = FileBytesRead(appdir + "gui/hidden.png")
	//hid = GTK_PixBuf_From_Bytes(bb, "png")
	//hid, _ = ResizePixelBuffer(hid, 64, 64)
}

type GtkFileIconBlock struct {
	ebox     *gtk.EventBox
	maingrid *gtk.Grid
	overlay  *gtk.Overlay
	overgrid *gtk.Grid

	icon         *gtk.Image
	icon_loading *gtk.Image
	icon_hidden  *gtk.Image
	icon_link    *gtk.Image
	icon_rw      *gtk.Image
	icon_mount   *gtk.Image
	check        *gtk.CheckButton
	name         *gtk.Label
	info         *gtk.Label

	fpath    string
	fname    string
	mime_sys string
	mime_app string
	app_open string
	isdir    bool
	isapp    bool
	islink   bool
	ishidden bool
	//rwx      int
}

func NewFileIconBlock(filepath string, filename string, wid int, isdir bool, islink bool, notread bool, ismount bool, strinfo string) *GtkFileIconBlock {
	b2 := BORDER_SIZE / 2
	isHidden := StringPart(filename, 1, 1) == "."

	icon, _ := gtk.ImageNew()
	icon.SetSizeRequest(ZOOM_SIZE, ZOOM_SIZE)

	check, _ := gtk.CheckButtonNew()
	/*check.SetHExpand(true)
	check.SetVExpand(true)*/
	check.SetHAlign(gtk.ALIGN_START)
	check.SetVAlign(gtk.ALIGN_START)
	check.SetMarginStart(b2)
	check.SetMarginTop(b2)

	icon_loading, _ := gtk.ImageNew()
	icon_loading.SetHAlign(gtk.ALIGN_CENTER)
	icon_loading.SetVAlign(gtk.ALIGN_CENTER)
	// icon_loading.SetHExpand(true)
	// icon_loading.SetVExpand(true)

	icon_link, _ := gtk.ImageNew()
	icon_link.SetHAlign(gtk.ALIGN_END)
	icon_link.SetVAlign(gtk.ALIGN_END)
	icon_link.SetMarginEnd(b2 * 3)
	icon_link.SetMarginBottom(b2 * 2)
	if islink {
		icon_link.SetFromPixbuf(pixbuf_link)
	}

	icon_rw, _ := gtk.ImageNew()
	icon_rw.SetHAlign(gtk.ALIGN_END)
	icon_rw.SetVAlign(gtk.ALIGN_START)
	icon_rw.SetMarginEnd(b2 * 3)
	icon_rw.SetMarginTop(b2 * 2)
	if notread {
		icon_rw.SetFromPixbuf(pixbuf_notread)
	}

	icon_mount, _ := gtk.ImageNew()
	icon_mount.SetHAlign(gtk.ALIGN_START)
	icon_mount.SetVAlign(gtk.ALIGN_END)
	icon_mount.SetMarginStart(b2 * 2)
	icon_mount.SetMarginBottom(b2 * 2)
	if ismount {
		icon_mount.SetFromPixbuf(pixbuf_mount)
	}

	overgrid, _ := gtk.GridNew()
	overgrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	overgrid.Attach(check, 0, 0, 1, 1)
	overgrid.Attach(icon_loading, 1, 1, 1, 1)
	overgrid.Attach(icon_rw, 2, 0, 1, 1)
	overgrid.Attach(icon_link, 2, 2, 1, 1)
	overgrid.Attach(icon_mount, 0, 2, 1, 1)
	overgrid.SetColumnHomogeneous(true)
	overgrid.SetRowHomogeneous(true)
	overgrid.SetSizeRequest(ZOOM_SIZE, ZOOM_SIZE)
	overgrid.SetHExpand(true)
	//overgrid.SetVExpand(true)

	icon_hidden, _ := gtk.ImageNew()
	if isHidden {
		icon_hidden.SetFromPixbuf(hid[ZOOM_SIZE])
	}

	overlay, _ := gtk.OverlayNew()
	overlay.Add(icon)
	overlay.AddOverlay(icon_hidden)
	overlay.AddOverlay(overgrid)
	overlay.SetSizeRequest(ZOOM_SIZE, ZOOM_SIZE)
	overlay.SetHExpand(true)
	//overlay.SetVExpand(true)

	name, _ := gtk.LabelNew(filename)
	GTK_LabelWrapMode(name, 3)
	name.SetJustify(gtk.JUSTIFY_CENTER)
	name.SetMarginStart(BORDER_SIZE)
	name.SetMarginEnd(BORDER_SIZE)
	name.SetHExpand(true)

	info, _ := gtk.LabelNew(strinfo)
	info.SetMarkup("<span color='#7F7F7F'>" + HtmlEscape(strinfo) + "</span>")
	GTK_LabelWrapMode(info, 2)
	info.SetJustify(gtk.JUSTIFY_CENTER)
	info.SetMarginStart(BORDER_SIZE)
	info.SetMarginEnd(BORDER_SIZE)
	info.SetHExpand(true)

	main, _ := gtk.GridNew()
	main.SetOrientation(gtk.ORIENTATION_VERTICAL)
	main.Attach(overlay, 0, 0, 1, 1)
	main.Attach(name, 0, 1, 1, 1)
	main.Attach(info, 0, 2, 1, 1)
	main.SetHExpand(true)

	evBox, _ := gtk.EventBoxNew()
	evBox.Add(main)
	// evBox.SetMarginStart(b2)
	// evBox.SetMarginEnd(b2)
	// evBox.SetMarginTop(b2)
	// evBox.SetMarginBottom(b2)
	evBox.Connect("draw", func(g *gtk.EventBox, ctx *cairo.Context) {
		if check.GetActive() {
			ctx.SetSourceRGBA(0.85, 0.9, 0.95, 1.0) // BLUE LIGHT
		} else {
			if !isHidden {
				//ctx.SetSourceRGBA(0, 0, 255, 1) //BACK_GRAY_VISIBLE
				ctx.SetSourceRGBA(BACK_GRAY_VISIBLE, BACK_GRAY_VISIBLE, BACK_GRAY_VISIBLE, 1)
			} else {
				ctx.SetSourceRGBA(BACK_GRAY_HIDDEN, BACK_GRAY_HIDDEN, BACK_GRAY_HIDDEN, 1)
			}
		}
		aw := g.GetAllocatedWidth()
		ah := g.GetAllocatedHeight()
		ctx.Rectangle(0, 0, float64(aw-2), float64(ah-2))
		/*for ry := 0; ry < ah-1; ry++ {
			for rx := 0; rx < aw-1; rx++ {
				if (rx+ry)%2 == 0 {
					ctx.Rectangle(float64(rx), float64(ry), 1, 1)
				}
			}
		}*/
		ctx.Fill()
		// tx2, ty2, _ := evBox.TranslateCoordinates(gGFiles, 0, 0)
		// Prln(filename + " [" + I2S(tx2) + "/" + I2S(ty2) + "]")
		// ctx.SetSourceRGBA(0.04, 0.07, 0.8, 1.0) // BLUE
		// ctx.Rectangle(float64(select_x1+tx2), float64(select_y1-0+ty2), float64(select_x2-select_x1), float64(select_y2-select_y1))
		// ctx.Fill()
	})
	check.Connect("button-release-event", func() {
		evBox.QueueDraw()
	})

	if isdir {
		dest := NewLinuxPath(true)
		dest.SetReal(filepath + filename)
		GTK_CopyPasteDnd_SetFolderDest(evBox, dest)
	}

	tfile := NewLinuxPath(isdir)
	tfile.SetReal(filepath + filename)
	getter := func() []*LinuxPath {
		list := []*LinuxPath{}
		fnames := FilesSelector_GetList()
		if len(fnames) > 1 {
			for j := 0; j < len(fnames); j++ {
				file1 := NewLinuxPath(false) //??
				file1.SetReal(filepath + fnames[j])
				list = append(list, file1)
			}
		} else {
			list = []*LinuxPath{tfile}
		}
		return list
	}
	GTK_CopyPasteDnd_SetIconSource(evBox, icon, getter)

	block := &GtkFileIconBlock{
		ebox:         evBox,
		maingrid:     main,
		overlay:      overlay,
		overgrid:     overgrid,
		icon:         icon,
		icon_loading: icon_loading,
		icon_link:    icon_link,
		icon_rw:      icon_rw,
		icon_mount:   icon_mount,
		check:        check,
		name:         name,
		info:         info,

		fname:    filename,
		fpath:    filepath,
		ishidden: isHidden,
		isdir:    isdir,
		islink:   islink,
	}
	block.SetWidth(wid)
	return block
}

func (i *GtkFileIconBlock) SetLoading(v bool) {
	if v {
		i.icon_loading.SetFromPixbuf(pixbuf_loading)
	} else {
		i.icon_loading.SetFromPixbuf(nil)
	}
}

func (i *GtkFileIconBlock) SetIconPixPuf(pixbuf_icon *gdk.Pixbuf) {
	i.icon.SetFromPixbuf(pixbuf_icon)
}

func (i *GtkFileIconBlock) ConnectEventBox(eventname string, f func(_ *gtk.EventBox, event *gdk.Event)) {
	i.ebox.Connect(eventname, f)
}

func (i *GtkFileIconBlock) SetWidth(wid int) {
	i.name.SetSizeRequest(wid, 32)
}

func (i *GtkFileIconBlock) GetFileName() string {
	return i.fname
}

func (i *GtkFileIconBlock) GetWidget() gtk.IWidget {
	return i.ebox
}

func (i *GtkFileIconBlock) IsClickedIn(x0, y0 int) bool {
	tx0, ty0, _ := i.ebox.TranslateCoordinates(gGFiles, 0, 0)
	tw := i.ebox.GetAllocatedWidth()
	th := i.ebox.GetAllocatedHeight()
	tx0 += BORDER_SIZE
	ty0 += BORDER_SIZE
	return x0 > tx0 && x0 < tx0+tw && y0 > ty0 && y0 < ty0+th
}

func (i *GtkFileIconBlock) IsInSelectRect(x1, y1, x2, y2 int) bool {
	tx0, ty0, _ := i.ebox.TranslateCoordinates(gGFiles, 0, 0)
	tw := i.ebox.GetAllocatedWidth()
	th := i.ebox.GetAllocatedHeight()
	tx0 += BORDER_SIZE
	ty0 += BORDER_SIZE
	checker := func(x0, y0, x1, y1, x2, y2 int) bool {
		return x0 > MINI(x1, x2) && x0 < MAXI(x1, x2) && y0 > MINI(y1, y2) && y0 < MAXI(y1, y2)
	}
	r1 := checker(tx0, ty0, x1, y1, x2, y2)
	r2 := checker(tx0+tw, ty0, x1, y1, x2, y2)
	r3 := checker(tx0, ty0+th, x1, y1, x2, y2)
	r4 := checker(tx0+tw, ty0+th, x1, y1, x2, y2)
	return r1 || r2 || r3 || r4
}

func (i *GtkFileIconBlock) SetSelected(v bool) {
	v0 := i.check.GetActive()
	if v0 != v {
		i.check.SetActive(v)
		i.ebox.QueueDraw()
	}
}

func (i *GtkFileIconBlock) Destroy() {
	i.icon.SetFromPixbuf(nil)

	//i.maingrid.Remove(i.icon)
	//i.maingrid.Remove(i.icon)
	//i.grid.Remove(i.check)
	//i.ebox.Remove(i.grid)

	i.icon.Destroy()
	i.icon_link.Destroy()
	i.icon_rw.Destroy()
	i.icon_mount.Destroy()
	i.check.Destroy()
	i.name.Destroy()
	i.info.Destroy()

	i.overgrid.Destroy()
	i.overlay.Destroy()
	i.maingrid.Destroy()
	//i.ebox.Destroy()
}
