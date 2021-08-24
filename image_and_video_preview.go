package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	. "github.com/SilentGopherLnx/GopherFileManager/pkg_fileicon"
	. "github.com/SilentGopherLnx/GopherFileManager/pkg_filetools"

	"github.com/gotk3/gotk3/gdk"

	"os"
	"os/exec"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"image"
	"image/color"
	"image/jpeg"

	"github.com/disintegration/imageorient"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/webp"
)

var MIME_IMAGE = []string{"jpg", "jpeg", "png", "gif", "webp", "ico", "bmp"}
var MIME_VIDEO = []string{"mp4", "avi", "mkv", "mov", "mpg", "mpeg", "flv", "wmv", "webm", "3gp"}
var MIME_PREVIEW = []string{}

var colorText color.RGBA

func init() {
	MIME_PREVIEW = append(MIME_IMAGE, MIME_VIDEO...)
	colorText = color.RGBA{200, 255, 255, 255}
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
func GetVideoPreviewBytes(filename string, zoom_size int, killchan chan *exec.Cmd) (*[]byte, int) {
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
	duration := S2I(seconds_str2)
	seconds := MINI(MAXI(1, duration*opt.GetVideoPercent()/100), 5000) //40*60=1800
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
	bb, _, _ := ExecCommandBytes([]byte{}, opt.GetFfmpegTimeout()*1000, killchan, "ffmpeg",
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
	return &bb, duration
}

func GetPreview_VideoPixBuf(filename string, zoom_size int, killchan chan *exec.Cmd, req int64, save_hash bool) (*gdk.Pixbuf, bool) {
	img2, ok2 := GetPreview_VideoImage(filename, zoom_size, killchan, req, save_hash)
	if ok2 {
		return GTK_PixBuf_From_RGBA(img2), true
	}
	return nil, false
}

func GetPreview_VideoImage(filename string, zoom_size int, killchan chan *exec.Cmd, req int64, save_hash bool) (image.Image, bool) {
	gr := uint8(RoundF(float64(255) * BACK_GRAY_VISIBLE))
	colorTransp := color.RGBA{gr, gr, gr, 0}

	zoom_max := Constant_ZoomMax()

	fbytes, dur := GetVideoPreviewBytes(filename, zoom_max, killchan)
	if len(*fbytes) == 0 {
		return nil, false
	}
	if req_id.Get() != req {
		Prln(">>>>>>>>> SKIP VIDEO: " + filename)
		return nil, false
	}
	img := ImageDecodeRGBA(fbytes, colorTransp)
	txt := "~" + I2S(RoundF(float64(dur)/60.0)) + "m"
	ImageText26x6_Bold(img, 5, 5, colorText, txt)
	if save_hash && !InterfaceNil(img) {
		info, err := FileInfo(filename, false)
		if err == nil {
			CachePreview_WriteImage(&info, 0, img, false)
		}
	}
	return GetPreview_ImageImage(img, zoom_size)
}

func ImageText26x6_Bold(img *image.RGBA, x, y int, col color.RGBA, label string) {
	colorBlack := color.RGBA{0, 0, 0, 255}
	for j := y - 1; j <= y+1; j++ {
		for i := x - 1; i <= x+1; i++ {
			if !(i == x && j == y) {
				ImageText26x6(img, i, j, colorBlack, label)
			}
		}
	}
	ImageText26x6(img, x, y, col, label)
}

func ImageText26x6(img *image.RGBA, x, y int, col color.RGBA, label string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6((y + 10) * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

/*func ImageText52x12(img *image.RGBA, x, y int, col color.RGBA, label string) {
	point := fixed.Point52_12{fixed.Int52_12(x * 64), fixed.Point52_12((y + 10) * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}*/

func GetPreview_ImagePixBuf(filename string, zoom_size int) (*gdk.Pixbuf, bool) {
	ftype := FileExtension(filename)
	if ftype == "jpg" {
		ftype = "jpeg"
	}
	fbytes, ok := FileBytesRead(filename)
	if ok {
		fr := imageorient.Decode
		if !opt.GetExifRot() {
			fr = nil
		}
		img := ImageDecodeCustom(fbytes, fr)
		img2, ok2 := GetPreview_ImageImage(img, zoom_size)
		if ok2 {
			return GTK_PixBuf_From_RGBA(img2), true
		}
	} else {
		Prln(filename + "//") //+ err.Error())
	}
	return nil, false
}

func GetPreview_ImageImage(img image.Image, zoom_size int) (*image.RGBA, bool) {
	if InterfaceNil(img) {
		return nil, false
	}
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	max_wh := MAXI(w, h)

	if max_wh <= zoom_size {
		img22 := image.NewRGBA(image.Rect(0, 0, zoom_size, zoom_size))
		ImageAddOver(img22, img, (zoom_size-w)/2, (zoom_size-h)/2)
		return img22, true
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
	return img3, true
}

// func EmptyIcon(w int, h int) *gdk.Pixbuf {
// 	pixbuf, _ := gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, w, h)
// 	return pixbuf
// }

func CachePreview_Function(info *FileReport, zoom_size int) string {
	hash_str := ""
	if info.IsDir() {
		hash_str = FilePathEndSlashRemove(info.FullName)
		//Prln("[[[[" + info.FullName + "]]]]")
	} else {
		hash_str = I2S64(info.Size()) + "/" + TimeStr(Time(info.ModTime()), true) + "/" + info.Name()
	}
	md5 := Crypto_MD5([]byte(hash_str))
	//md5 = StringReplace(hash_str, "/", "_")
	//Prln("md5:" + md5)
	//Prln("sha1:" + Crypto_SHA1([]byte(hash_str)))
	if info.IsDir() {
		return "D-" + I2S(zoom_size) + "_" + md5
	} else {

		return "F-" + StringUp(FileExtension(info.Name())) + "_" + md5
	}
}

func CachePreview_ReadPixbuf(info *FileReport, zoom_size int, alphamask *image.RGBA) *gdk.Pixbuf {
	return GTK_PixBuf_From_RGBA(CachePreview_ReadImage(info, zoom_size, alphamask))
}

func CachePreview_ReadImage(info *FileReport, zoom_size int, alphamask *image.RGBA) *image.RGBA {
	hash_str := CachePreview_Function(info, zoom_size)
	data, ok := FileBytesRead(opt.GetHashFolder() + hash_str + ".jpg")

	var img *image.RGBA
	img = nil
	if ok {
		if alphamask == nil {

			//return GTK_PixBuf_From_Bytes(data, "jpeg")

			gr := uint8(RoundF(float64(255) * BACK_GRAY_VISIBLE))
			colorT := color.RGBA{gr, gr, gr, 0}
			img = ImageDecodeRGBA(data, colorT)

			// nomask := ImageDecode(data)
			// img := image.NewRGBA(image.Rect(0, 0, zoom_size, zoom_size))
			// ImageAddOver(img, nomask, 0, 0)
			// return img
		} else {
			nomask := ImageDecode(data)
			img = image.NewRGBA(image.Rect(0, 0, zoom_size, zoom_size))
			ImageAddOver(img, nomask, 0, 0)
			for y := 0; y < zoom_size; y++ {
				for x := 0; x < zoom_size; x++ {
					col_nomask := img.RGBAAt(x, y)
					col_amask := alphamask.RGBAAt(x, y)
					col_nomask.A = col_amask.A
					img.SetRGBA(x, y, col_nomask)
				}
			}
		}
	}
	if img != nil {
		if !info.IsDir() && img.Rect.Max.X != zoom_size {
			//return ImageResizeNearest(img, zoom_size, zoom_size)
			img2, ok2 := GetPreview_ImageImage(img, zoom_size)
			if ok2 {
				return img2
			}
		}
		return img
	} else {
		return nil
	}
}

func CachePreview_WriteImage(info *FileReport, zoom_size int, img image.Image, delete_mode bool) bool {
	hash_str := CachePreview_Function(info, zoom_size)
	cname := opt.GetHashFolder() + hash_str + ".jpg"
	if delete_mode {
		FileDelete(cname)
		return true
	} else {
		f, err1 := os.Create(cname)
		if err1 == nil {
			err2 := jpeg.Encode(f, img, &jpeg.Options{Quality: 50})
			if err2 == nil {
				return true
			}
		}
		return false
	}
}

func GetPixBufGTK_Folder(folderpath string, zoom_size int, basic_mode bool, qu *SyncQueue, killchan chan *exec.Cmd, req int64, skip_cached bool) (*gdk.Pixbuf, bool) { // icon_msg *IconUpdateable
	folderpath2 := FolderPathEndSlash(folderpath)
	imgRGBA := GetIcon_ImageFolder(zoom_size)
	scale := zoom_size / 64
	xy := [][]int{
		[]int{2 * scale, 17 * scale},
		[]int{28 * scale, 17 * scale},
		[]int{12 * scale, 31 * scale},
		[]int{38 * scale, 31 * scale},
	}
	imgs := []image.Image{}
	total_icons := 4 //len(xy)

	numdirs := 0
	numfiles := 0
	files, err := Folder_ListFiles(folderpath, false)
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
				imgs2, changed = replace_mime_images(folderpath2, files, imgs2, mimetypes, zoom_size, killchan, req, skip_cached)
			}
			mimetypes[0] = mimetypes[0] //used for compile
		}
		imgs = append(imgs1, imgs2...)

		imlen := MINI(len(imgs), total_icons)
		for j := 0; j < imlen; j++ {
			ImageAddOver(imgRGBA, imgs[j], xy[j][0], xy[j][1])
			if j == 0 && numdirs > icondirs && icondirs == 1 {
				ImageAddOver(imgRGBA, imgs[0], xy[j][0]+scale*2, xy[j][1]+scale*3)
			}
		}

		if numdirs > 2 {
			ImageText26x6_Bold(imgRGBA, zoom_size/8, zoom_size*2/5, colorText, I2S(numdirs))
			//ImageText52x12(imgRGBA, zoom_size/8, zoom_size*2/5, colorText, I2S(numdirs))
		}
		if numfiles > 3 {
			ImageText26x6_Bold(imgRGBA, zoom_size*6/8, zoom_size*7/10, colorText, I2S(numfiles))
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
		if req_id.Get() != req {
			Prln(">>>>>>>>> SKIP FOLDER CACHE PIC: " + folderpath2)
		} else {
			CachePreview_WriteImage(&FileReport{FullName: FilePathEndSlashRemove(folderpath2), IsDirectory: true}, zoom_size, imgRGBA, false)
		}
	} else {
		CachePreview_WriteImage(&FileReport{FullName: FilePathEndSlashRemove(folderpath2), IsDirectory: true}, zoom_size, imgRGBA, true)
	}

	pixbuf := GTK_PixBuf_From_RGBA(imgRGBA)
	if pixbuf != nil {
		return pixbuf, true
	}
	return nil, false
}

func best_file_icons(folderpath string, files []FileReport, maxicons int, zoom_size int, basic_mode bool) ([]image.Image, []string) {
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

func replace_mime_images(folderpath string, files []FileReport, imgs []image.Image, mimes []string, zoom_size int, killchan chan *exec.Cmd, req int64, skip_cached bool) ([]image.Image, bool) {
	mime_image := PREFIX_DRAWONME + FILE_TYPE_IMAGE
	mime_video := PREFIX_DRAWONME + FILE_TYPE_MOVIE
	imgs2 := []image.Image{}
	file_j := 0
	changed := false
	fr := imageorient.Decode
	if !opt.GetExifRot() {
		fr = nil
	}
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
						if req_id.Get() != req {
							Prln("////////////////// SKIP: " + f.Name())
						} else {
							tfile := FileExtension(f.Name())
							if StringInArray(tfile, MIME_IMAGE) > -1 {
								fbytes, ok := FileBytesRead(folderpath + f.Name())
								if ok {
									img_new := ImageDecodeCustom(fbytes, fr)
									if img_new != nil {
										zoom_size2 := ZoomSmall(zoom_size)
										zoom_size3 := zoom_size2 - zoom_size/32
										w := img_new.Bounds().Max.X
										h := img_new.Bounds().Max.Y
										if w > 1 && h > 1 {
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
										}
									} else {
										Prln("??" + f.Name())
									}
								}
							}
							if StringInArray(tfile, MIME_VIDEO) > -1 {
								zoom_size2 := ZoomSmall(zoom_size)
								skip_cached_ok := false
								if skip_cached {
									img_new := CachePreview_ReadImage(&f, zoom_size2, nil)
									if img_new != nil {
										exist = true
										imgs2 = append(imgs2, img_new)
										skip_cached_ok = true
										Prln("not recashing: " + f.FullName)
									}
								}
								if !skip_cached_ok {
									img_new, ok2 := GetPreview_VideoImage(f.FullName, zoom_size2, killchan, req, true)
									if ok2 {
										exist = true
										imgs2 = append(imgs2, img_new)
									}
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
