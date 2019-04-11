package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"image/color"
	"image/jpeg"

	"image"
	"os"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/webp"
)

var MIME_IMAGE = []string{"jpg", "jpeg", "png", "gif", "webp", "ico", "bmp"}
var MIME_VIDEO = []string{"mp4", "avi", "mkv", "mov", "mpg", "mpeg", "flv", "wmv", "webm", "3gp"}
var MIME_PREVIEW = []string{}

func init() {
	MIME_PREVIEW = append(MIME_IMAGE, MIME_VIDEO...)
}

func FileIsPreviewAbble(tfile string) bool {
	return StringInArray(tfile, MIME_PREVIEW) > -1
}

func getConfigValue(allinfo []string, name string, skip string) string {
	skip_down := StringDown(skip)
	for j := 0; j < len(allinfo); j++ {
		if StringFind(allinfo[j], name) == 1 {
			value := StringTrim(StringPart(allinfo[j], StringLength(name)+1, 0))
			if StringDown(value) != skip_down {
				return value
			}
		}
	}
	return ""
}

// sudo apt install ffmpeg
// sudo apt install ffprobe
func GetVideoPreviewBytes(filename string, zoom_size int) *[]byte {
	info, _, _ := ExecCommand("ffprobe", "-i", filename, "-show_format", "-show_streams")
	//Prln("[" + filename + "]:" + info)
	info_arr := StringSplitLines(info)
	w_old := S2I(getConfigValue(info_arr, "width=", "0"))
	h_old := S2I(getConfigValue(info_arr, "height=", "0"))
	seconds_str1 := getConfigValue(info_arr, "duration=", "N/A")
	seconds_str2 := seconds_str1
	ind := StringFind(seconds_str2, ".")
	if ind > 0 {
		seconds_str2 = StringPart(seconds_str2, 1, ind-1)
	}
	seconds := MINI(MAXI(1, S2I(seconds_str2)/2), 2400) //40*60=1800
	if w_old == 0 || h_old == 0 {
		w_old = 16
		h_old = 9
	}
	max_old := MAXI(w_old, h_old)
	max_new := zoom_size - 4
	w_new := MAXI(1, max_new*w_old/max_old)
	h_new := MAXI(1, max_new*h_old/max_old)

	ss := I2S(seconds/60) + ":" + StringEnd("0"+I2S(seconds%60), 2) + ".0"
	wh := I2S(w_new) + "x" + I2S(h_new)
	//Prln("[" + wh + "/" + I2S(w_old) + ":" + I2S(h_old) + "] [" + ss + "/" + seconds_str1 + "] " + filename)
	//"-r", "10"
	//"-preset", "ultrafast"
	//"-vcodec", "libx264"

	// -ss before -i is quicker !!!!!!
	bb, _, _ := ExecCommandBytes([]byte{}, opt.GetFfmpegTimeout()*1000, "ffmpeg",
		"-ss", ss,
		"-i", filename,
		"-vframes", "1",
		"-s", wh,
		"-preset", "ultrafast",
		"-tune", "fastdecode",
		"-crf", "80",
		//"-vcodec", "libx264",
		"-f", "singlejpeg", "-")

	//Prln("BYTES:" + string(bb))
	return &bb
}

func GetPixBufGTK_Video(filename string, zoom_size int, save_hash bool) (*gdk.Pixbuf, bool) {
	// width := zoom_size - 4
	// height := width * 9 / 16
	gr := uint8(RoundF(float64(255) * BACK_GRAY_VISIBLE))
	colorT := color.RGBA{gr, gr, gr, 0}
	ffmpeg := true
	if ffmpeg {
		fbytes := GetVideoPreviewBytes(filename, zoom_size)
		if len(*fbytes) == 0 {
			return nil, false
		}
		img := ImageDecode(fbytes)
		img2 := image.NewRGBA(image.Rect(0, 0, zoom_size, zoom_size))
		w_ := img.Bounds().Max.X
		h_ := img.Bounds().Max.Y
		for y := 0; y < zoom_size; y++ {
			for x := 0; x < zoom_size; x++ {
				img2.SetRGBA(x, y, colorT)
			}
		}
		ImageAddOver(img2, img, (zoom_size-w_)/2, (zoom_size-h_)/2)
		pixbuf := GTK_PixBuf_From_RGBA(img2)
		if pixbuf == nil {
			return nil, false
		} else {
			if save_hash {
				WriteHashImage(filename, zoom_size, img2)
			}
			return pixbuf, true
		}
	} else {
		// img := GetVideoFrame_Mpeg2(filename)
		// if img != nil {
		// 	img2 := ImageResizeNearest(img, width*2, height*2)
		// 	img3 := ImageResizeHalfNice(img2)
		// 	if img3 != nil {
		// 		pixbuf := GTK_PixBuf_From_RGBA(img)
		// 		if pixbuf != nil {
		// 			return pixbuf, true
		// 		}
		// 	}
		// }
		return nil, false
	}
}

// func GetPixBufGTK_Image_v01(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
// 	//pixbuf, err := gdk.PixbufNewFromFile(filename)
// 	ftype := FileExtension(filename)
// 	if ftype == "jpg" {
// 		ftype = "jpeg"
// 	}
// 	fbytes, ok := FileBytesRead(filename)
// 	if ok {
// 		img := ImageDecode(fbytes)
// 		img2 := ImageResizeNearest(img, zoom_size, zoom_size)
// 		img = nil
// 		fbytes = nil

// 		/*data := []byte{}
// 		buf := new(bytes.Buffer)
// 		err := jpeg.Encode(buf, img2, &jpeg.Options{Quality: 90})
// 		if err == nil {
// 			data = buf.Bytes()
// 		}
// 		img = nil
// 		img2 = nil
// 		data2 := CloneBytesArray(data)
// 		return GTK_PixBuf_From_Bytes(&data2, "jpeg"), true*/

// 		//return nil, false
// 		pixbuf := GTK_PixBuf_From_RGBA(img2)
// 		return pixbuf, true
// 	} else {
// 		Prln(filename + "//") //+ err.Error())
// 	}
// 	return nil, false
// }

func GetPixBufGTK_Image(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
	//pixbuf, err := gdk.PixbufNewFromFile(filename)
	ftype := FileExtension(filename)
	if ftype == "jpg" {
		ftype = "jpeg"
	}
	fbytes, ok := FileBytesRead(filename)
	if ok {
		img := ImageDecode(fbytes)
		if InterfaceNil(img) {
			return nil, false
		}
		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y
		max_wh := MAXI(w, h)

		if max_wh <= zoom_size {
			img22 := image.NewRGBA(image.Rect(0, 0, zoom_size, zoom_size))
			ImageAddOver(img22, img, (zoom_size-w)/2, (zoom_size-h)/2)
			return GTK_PixBuf_From_RGBA(img22), true
		}

		size_up := 2
		if max_wh >= zoom_size*4 {
			size_up = 4
		}
		max_new := (zoom_size - 4) * size_up
		zoom_size_up := zoom_size * size_up

		w2 := MAXI(1, max_new*w/max_wh)
		h2 := MAXI(1, max_new*h/max_wh)
		img2 := ImageResizeNearest(img, w2, h2)
		img22 := image.NewRGBA(image.Rect(0, 0, zoom_size_up, zoom_size_up))
		ImageAddOver(img22, img2, (zoom_size_up-w2)/2, (zoom_size_up-h2)/2)
		img3 := ImageResizeHalfNice(img22)
		if size_up == 4 {
			img3 = ImageResizeHalfNice(img3)
		}
		//img3 := ImageResizeNearest(img22, zoom_size, zoom_size)
		return GTK_PixBuf_From_RGBA(img3), true
	} else {
		Prln(filename + "//") //+ err.Error())
	}
	return nil, false
}

// func GetPixBufGTK_Image(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
// 	return GetPixBufGTK_Image_v00(filename, zoom_size)
// }

/*func GetPixBufGTK_Image_v1(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
	pixbuf, err := gdk.PixbufNewFromFile(filename)
	if err == nil {
		max_wh := MAXI(pixbuf.GetWidth(), pixbuf.GetHeight())
		pixbuf2 := pixbuf
		ok := false
		if max_wh > zoom_size*2 {
			pixbuf2, ok = ResizePixelBuffer(pixbuf, zoom_size*2, gdk.INTERP_NEAREST)
			if !ok {
				return nil, false
			}
		}
		if max_wh < zoom_size {
			return pixbuf, true
		}
		//pixbuf.Ref()
		pixbuf3, ok3 := ResizePixelBuffer(pixbuf2, zoom_size, gdk.INTERP_BILINEAR)
		//pixbuf2.Unref()
		return pixbuf3, ok3
	} else {
		Prln(filename + "//") //+ err.Error())
	}
	return nil, false
}

func GetPixBufGTK_Image_v2(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
	ftype := GetFileExtension(filename)
	if ftype == "jpg" {
		ftype = "jpeg"
	}
	fbytes, ok := FileBytesRead(filename)
	if ok {
		pixbuf := GTK_PixBuf_From_Bytes(fbytes, ftype)
		//if err == nil {
		max_wh := MAXI(pixbuf.GetWidth(), pixbuf.GetHeight())
		pixbuf2 := pixbuf
		ok := false
		if max_wh > zoom_size*2 {
			pixbuf2, ok = ResizePixelBuffer(pixbuf, zoom_size*2, gdk.INTERP_NEAREST)
			if !ok {
				return nil, false
			}
		}
		if max_wh < zoom_size {
			return pixbuf, true
		}
		//pixbuf.Ref()
		pixbuf3, ok3 := ResizePixelBuffer(pixbuf2, zoom_size, gdk.INTERP_BILINEAR)
		//pixbuf2.Unref()
		return pixbuf3, ok3
	} else {
		Prln(filename + "//") //+ err.Error())
	}
	return nil, false
}*/

func ResizePixelBuffer(pixbuf *gdk.Pixbuf, zoom_size int, interp gdk.InterpType) (*gdk.Pixbuf, bool) {
	w_old := pixbuf.GetWidth()
	h_old := pixbuf.GetHeight()
	max_old := MAXI(w_old, h_old)
	max_new := zoom_size - 4
	w_new := MAXI(1, max_new*w_old/max_old)
	h_new := MAXI(1, max_new*h_old/max_old)
	pixbuf, err := pixbuf.ScaleSimple(w_new, h_new, interp)
	if err == nil {
		return pixbuf, true
	}
	return nil, false
}

// func EmptyIcon(w int, h int) *gdk.Pixbuf {
// 	pixbuf, _ := gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, w, h)
// 	return pixbuf
// }

func ReadHashPixbuf(path string, zoom_size int) *gdk.Pixbuf {
	md5 := Crypto_MD5([]byte(I2S(zoom_size) + "//" + path))
	//Prln("md5:" + md5)
	//Prln("sha1:" + Crypto_SHA1([]byte(path)))
	data, ok := FileBytesRead(opt.GetHashFolder() + md5 + ".jpg")
	if ok {
		return GTK_PixBuf_From_Bytes(data, "jpeg")
	} else {
		return nil
	}
}

func WriteHashImage(path string, zoom_size int, img image.Image) bool {
	md5 := Crypto_MD5([]byte(I2S(zoom_size) + "//" + path))
	//Prln("md5:" + md5)
	//Prln("sha1:" + Crypto_SHA1([]byte(path)))
	f, err1 := os.Create(opt.GetHashFolder() + md5 + ".jpg")
	if err1 == nil {
		err2 := jpeg.Encode(f, img, &jpeg.Options{Quality: 50})
		if err2 == nil {
			return true
		}
	}
	return false
}

func GetPixBufGTK_Folder(folderpath string, zoom_size int, basic_mode bool, qu *SyncQueue, icon_msg *IconUpdateable) (*gdk.Pixbuf, bool) {
	folderpath2 := FolderPathEndSlash(folderpath)
	imgRGBA := GetIcon_ImageFolder(zoom_size)
	scale := zoom_size / 64
	xy := [][]int{
		[]int{2 * scale, 17 * scale},
		[]int{26 * scale, 17 * scale},
		[]int{14 * scale, 31 * scale},
		[]int{38 * scale, 31 * scale},
	}
	imgs := []image.Image{}
	total_icons := 4 //len(xy)

	numdirs := 0
	numfiles := 0
	files, err := Folder_ListFiles(folderpath)
	if err == nil {
		for _, f := range files {
			if f.IsDir() {
				numdirs++
			} else {
				numfiles++
			}
		}

		icondirs := 0
		iconfiles := 0

		if numdirs > 0 {
			icondirs = MINI(numdirs, MAXI(1, total_icons-numfiles))
			iconfiles = MINI(total_icons-icondirs, numfiles)
		} else {
			iconfiles = MINI(total_icons, numfiles)
		}

		imgs1 := []image.Image{}
		if icondirs > 0 {
			dicon := GetIcon_ImageRGBA(zoom_size, "", true)
			for j := 0; j < icondirs; j++ {
				imgs1 = append(imgs1, dicon)
			}
		}

		var imgs2 []image.Image
		var mimetypes []string
		changed := false
		if iconfiles > 0 {
			imgs2, mimetypes = best_file_icons(folderpath2, files, iconfiles, zoom_size, basic_mode)
			if !basic_mode {
				imgs2, changed = replace_mime_images(folderpath2, files, imgs2, mimetypes, zoom_size)
			}
			mimetypes[0] = mimetypes[0] //used for compile
		}
		imgs = append(imgs1, imgs2...)

		imlen := MINI(len(imgs), total_icons)
		for j := 0; j < imlen; j++ {
			ImageAddOver(imgRGBA, imgs[j], xy[j][0], xy[j][1])
		}

		changed = changed

		/*go func() {
			if iconfiles > 0 && !basic_mode {
				imgs2, changed = replace_mime_images(folderpath2, files, imgs2, mimetypes, zoom_size)
				if changed {
					imgs = append(imgs1, imgs2...)
					imgRGBA := GetIcon_ImageFolder(zoom_size)
					imlen := MINI(len(imgs), total_icons)
					for j := 0; j < imlen; j++ {
						ImageAddOver(imgRGBA, imgs[j], xy[j][0], xy[j][1])
					}
					qu.Append(IconUpdateable{icon: msg.icon, fullname: msg.fullname, fname: msg.fname, tfile: msg.tfile, basic_mode: false, folder: true, pixbuf_preview})
				}
			}
		}()*/
	}

	if numfiles > 0 || numdirs > 0 {
		WriteHashImage(folderpath2, zoom_size, imgRGBA)
	}

	pixbuf := GTK_PixBuf_From_RGBA(imgRGBA)
	if pixbuf != nil {
		return pixbuf, true
	}
	return nil, false
}

func best_file_icons(folderpath string, files []os.FileInfo, maxicons int, zoom_size int, basic_mode bool) ([]image.Image, []string) {
	type pair struct {
		mime  string
		count int
	}
	pairs := []*pair{}
	numfiles := 0
	for _, f := range files {
		if !f.IsDir() {
			numfiles++

			mime_new := ""
			mime_system := ""
			if !basic_mode {
				mime_system = FileMIME(folderpath + f.Name())
			}
			if mime_system == APP_EXEC_TYPE {
				mime_new = PREFIX_EXTRA + FILE_TYPE_BIN
			} else {
				tfile := FileExtension(f.Name())
				mime_new = GetExtensionIconName(tfile, false)
				if f.Size() == 0 {
					if f.Mode().IsRegular() {
						mime_new = PREFIX_DRAWONME + FILE_TYPE_ZERO
					} else {
						mime_new = PREFIX_DRAWONME + FILE_TYPE_NOTFILE
					}
					//Prln("zero")
				}
			}

			exist := false
			for _, p := range pairs {
				if mime_new == p.mime {
					p.count++
					exist = true
				}
			}
			if !exist {
				pairs = append(pairs, &pair{mime_new, 1})
			}
		}
	}

	sort_pairs := func() {
		SortArray(pairs, func(i, j int) bool {
			return pairs[i].count > pairs[j].count
		})
	}
	minlen := MINI(maxicons, len(pairs))
	maxlen := MINI(maxicons, numfiles)

	sort_pairs()

	if minlen < maxicons && numfiles > minlen {
		for len(pairs) < maxlen {
			pairs[0].count--
			pairs = append(pairs, pairs[0])
			sort_pairs()
		}
		minlen = MINI(maxicons, len(pairs))
	}

	imgs := []image.Image{}
	mimes := []string{}
	for j := 0; j < minlen; j++ {
		imgs = append(imgs, GetIcon_ImageRGBA_OF(zoom_size, pairs[j].mime))
		mimes = append(mimes, pairs[j].mime)
	}
	return imgs, mimes
}

func replace_mime_images(folderpath string, files []os.FileInfo, imgs []image.Image, mimes []string, zoom_size int) ([]image.Image, bool) {
	mime_image := PREFIX_DRAWONME + FILE_TYPE_IMAGE
	mime_video := PREFIX_DRAWONME + FILE_TYPE_MOVIE
	imgs2 := []image.Image{}
	file_j := 0
	changed := false
	for k := 0; k < len(imgs); k++ {
		if mimes[k] != mime_image && mimes[k] != mime_video {
			imgs2 = append(imgs2, imgs[k])
		} else {
			exist := false
			for j := file_j; j < len(files); j++ {
				if !exist {
					file_j = j + 1
					f := files[j]
					if !f.IsDir() {
						tfile := FileExtension(f.Name())
						if StringInArray(tfile, MIME_IMAGE) > -1 {
							fbytes, ok := FileBytesRead(folderpath + f.Name())
							if ok {
								img_new := ImageDecode(fbytes)
								if img_new != nil {
									zoom_size2 := ZoomSmall(zoom_size)
									zoom_size3 := zoom_size2 - zoom_size/32
									w := img_new.Bounds().Max.X
									h := img_new.Bounds().Max.Y
									max_old := MAXI(w, h)
									w_new := w
									h_new := h
									var img_new2 *image.RGBA
									if max_old > zoom_size3 {
										w_new = MAXI(1, zoom_size3*w*2/max_old)
										h_new = MAXI(1, zoom_size3*h*2/max_old)
										img_new = ImageResizeNearest(img_new, w_new, h_new)
										img_new2 = image.NewRGBA(image.Rect(0, 0, zoom_size2*2, zoom_size2*2))
										ImageAddOver(img_new2, img_new, zoom_size2-w_new/2, zoom_size2-h_new/2)
										img_new2 = ImageResizeHalfNice(img_new2)
									} else {
										img_new2 = image.NewRGBA(image.Rect(0, 0, zoom_size2, zoom_size2))
										ImageAddOver(img_new2, img_new, (zoom_size2-w_new)/2, (zoom_size2-h_new)/2)
									}
									exist = true
									imgs2 = append(imgs2, img_new2)
								} else {
									Prln("??" + f.Name())
								}
							}
						}
						if StringInArray(tfile, MIME_VIDEO) > -1 {
							zoom_size2 := ZoomSmall(zoom_size)
							// w_new := zoom_size2 - 2
							// h_new := w_new * 9 / 16
							fbytes := GetVideoPreviewBytes(folderpath+f.Name(), zoom_size2)
							if fbytes != nil && len(*fbytes) > 0 {
								img_new := ImageDecode(fbytes)
								if img_new != nil {
									img_new2 := image.NewRGBA(image.Rect(0, 0, zoom_size2, zoom_size2))
									w_new := img_new.Bounds().Max.X
									h_new := img_new.Bounds().Max.Y
									ImageAddOver(img_new2, img_new, (zoom_size2-w_new)/2, (zoom_size2-h_new)/2)
									exist = true
									imgs2 = append(imgs2, img_new2)
								}
							}
						}
					}
				}
			}
			if exist {
				changed = true
			} else {
				imgs2 = append(imgs2, imgs[k])
			}
		}
	}
	return imgs2, changed
}

type IconUpdateable struct {
	icon           *gtk.Image
	loading        *gtk.Image
	fullname       string
	fname          string
	tfile          string
	pixbuf_preview *gdk.Pixbuf
	basic_mode     bool
	folder         bool
	oldbuf         bool //have loaded old preview
}

func IconThread(icon_chan chan *IconUpdateable, qu *SyncQueue, thread_id int) { // qu *queue.Queue
	for {
		/*Prln("[" + I2S(thread_id) + "]Waiting..")
		runtime.Gosched()
		Sleep(5)*/
		msg := <-icon_chan
		num_works.Add(1)
		if GTK_WidgetExist(msg.icon) {
			//var pixbuf_preview *gdk.Pixbuf
			var ok = false
			if msg.folder {
				//if !msg.basic_mode {
				msg.pixbuf_preview, ok = GetPixBufGTK_Folder(msg.fullname, ZOOM_SIZE, msg.basic_mode, qu, msg)
				//}
			} else {
				if StringInArray(msg.tfile, MIME_IMAGE) > -1 {
					msg.pixbuf_preview, ok = GetPixBufGTK_Image(msg.fullname, ZOOM_SIZE)
				}
				if !msg.basic_mode && StringInArray(msg.tfile, MIME_VIDEO) > -1 {
					msg.pixbuf_preview, ok = GetPixBufGTK_Video(msg.fullname, ZOOM_SIZE, true)
				}
			}
			if ok {
				qu.Append(msg)
			}
		}
		num_works.Add(-1)
		RuntimeGosched()
	}
}
