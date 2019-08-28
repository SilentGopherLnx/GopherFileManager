package pkg_fileicon

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	"image"
	"image/color"

	"github.com/gotk3/gotk3/gdk"
)

const MIME_PATH = "mime/"
const MIME_CONFIG = "mime/filetypes.cfg"

const PREFIX_DRAWONME = "drawonme_"
const PREFIX_EXTRA = "extra_"

const FILE_TYPE_FOLDER = "folder"
const FILE_TYPE_FOLDER_HASH = "folder_hash"
const FILE_TYPE_UNKNOWN = "unknown"
const FILE_TYPE_ZERO = "zero"
const FILE_TYPE_NOTFILE = "notfile"
const FILE_TYPE_IMAGE = "image"
const FILE_TYPE_MOVIE = "movie"
const FILE_TYPE_ARCHIVE = "archive"
const FILE_TYPE_BIN = "bin"

var TABLE_EXTENSIONS_ICONS_NAMES map[string]string
var TABLE_EXTENSIONS_ICONS_SETS map[int]*IconSetOfZoom

func init() {
	appdir := FolderLocation_App()

	TABLE_EXTENSIONS_ICONS_NAMES = make(map[string]string)
	txt, ok := FileTextRead(appdir + MIME_CONFIG)
	if ok {
		strs := StringSplitLines(txt)
		prefix := ""
		for _, str := range strs {
			if len(str) > 0 {
				if StringPart(str, 1, 1) == "#" {
					prefix = StringPart(str, 2, 0)
				} else {
					ab := StringSplit(str, " ")
					if len(ab) >= 2 {
						types := StringSplit(ab[1], ",")
						for _, ftype := range types {
							TABLE_EXTENSIONS_ICONS_NAMES[ftype] = prefix + "_" + ab[0]
							//Prln(ftype + "/" + prefix + "/" + ab[0])
						}
					}
				}
			}
		}
	} else {
		TABLE_EXTENSIONS_ICONS_NAMES["jpeg"] = PREFIX_DRAWONME + FILE_TYPE_IMAGE
		TABLE_EXTENSIONS_ICONS_NAMES["jpg"] = PREFIX_DRAWONME + FILE_TYPE_IMAGE
		TABLE_EXTENSIONS_ICONS_NAMES["png"] = PREFIX_DRAWONME + FILE_TYPE_IMAGE
		TABLE_EXTENSIONS_ICONS_NAMES["gif"] = PREFIX_DRAWONME + FILE_TYPE_IMAGE

		TABLE_EXTENSIONS_ICONS_NAMES["mp4"] = PREFIX_DRAWONME + FILE_TYPE_MOVIE
		TABLE_EXTENSIONS_ICONS_NAMES["avi"] = PREFIX_DRAWONME + FILE_TYPE_MOVIE
		TABLE_EXTENSIONS_ICONS_NAMES["mkv"] = PREFIX_DRAWONME + FILE_TYPE_MOVIE
		TABLE_EXTENSIONS_ICONS_NAMES["mov"] = PREFIX_DRAWONME + FILE_TYPE_MOVIE
		TABLE_EXTENSIONS_ICONS_NAMES["mpg"] = PREFIX_DRAWONME + FILE_TYPE_MOVIE

		TABLE_EXTENSIONS_ICONS_NAMES["zip"] = PREFIX_DRAWONME + FILE_TYPE_ARCHIVE
	}

	zooms := Constant_ZoomArray()
	TABLE_EXTENSIONS_ICONS_SETS = make(map[int]*IconSetOfZoom)
	for _, zoom := range zooms {
		TABLE_EXTENSIONS_ICONS_SETS[zoom] = Make_IconSetOfZoom(zoom)
	}
}

func Constant_ZoomArray() []int {
	return []int{64, 128, 256} //, 512}
}

func Constant_ZoomMax() int {
	arr := Constant_ZoomArray()
	return arr[len(arr)-1]
}

func ZoomSmall(zoom_big int) int {
	return RoundF(float64(zoom_big) / 8.0 * 3.0)
}

type IconSetOfZoom struct {
	BigPixBufs      map[string]*gdk.Pixbuf
	SmallImageRGBA  map[string]image.Image
	FolderImageRGBA *image.RGBA
}

func Make_IconSetOfZoom(zoom int) *IconSetOfZoom {
	gr := uint8(RoundF(float64(255) * BACK_GRAY_VISIBLE))
	colorT := color.RGBA{gr, gr, gr, 0}
	zoom_small := ZoomSmall(zoom)
	iconset := &IconSetOfZoom{BigPixBufs: make(map[string]*gdk.Pixbuf), SmallImageRGBA: make(map[string]image.Image)}

	folder_big, _ := icon_set_big_small(iconset, zoom, zoom_small, PREFIX_DRAWONME+FILE_TYPE_FOLDER, &colorT)
	iconset.FolderImageRGBA = ImageDecodeRGBA(folder_big, colorT)
	if InterfaceNil(iconset.FolderImageRGBA) {
		iconset.FolderImageRGBA = image.NewRGBA(image.Rect(0, 0, zoom, zoom))
	}

	icon_set_big_small(iconset, zoom, zoom_small, PREFIX_DRAWONME+FILE_TYPE_FOLDER_HASH, &colorT)
	icon_set_big_small(iconset, zoom, zoom_small, PREFIX_DRAWONME+FILE_TYPE_UNKNOWN, &colorT)
	icon_set_big_small(iconset, zoom, zoom_small, PREFIX_DRAWONME+FILE_TYPE_ZERO, &colorT)
	icon_set_big_small(iconset, zoom, zoom_small, PREFIX_DRAWONME+FILE_TYPE_NOTFILE, &colorT)

	for _, val := range TABLE_EXTENSIONS_ICONS_NAMES {
		_, ok := iconset.BigPixBufs[val]
		if !ok {
			icon_set_big_small(iconset, zoom, zoom_small, val, &colorT)
		}
	}

	return iconset
}

func icon_set_big_small(iconset *IconSetOfZoom, zoom int, zoom_small int, name string, colorT *color.RGBA) (*[]byte, *[]byte) {
	appdir := FolderLocation_App()
	bytes_big := try_load_image_of_zoom(appdir+MIME_PATH, 0, zoom, "/"+name+".png")
	bytes_small := try_load_image_of_zoom(appdir+MIME_PATH, 1, zoom_small, "/"+name+".png")

	pixbuf := GTK_PixBuf_From_Bytes(bytes_big, "png")
	if pixbuf != nil && (pixbuf.GetWidth() != zoom || pixbuf.GetHeight() != zoom) {
		pixbuf, _ = ResizePixelBuffer(pixbuf, zoom, gdk.INTERP_BILINEAR)
	}

	image := ImageDecode(bytes_small)
	if pixbuf != nil && (pixbuf.GetWidth() != zoom_small || pixbuf.GetHeight() != zoom_small) {
		image2 := ImageResizeNearest(image, zoom_small*2, zoom_small*2)
		image = ImageResizeHalfNice(image2)
	}

	iconset.BigPixBufs[name] = pixbuf
	iconset.SmallImageRGBA[name] = image
	return bytes_big, bytes_small
}

func try_load_image_of_zoom(prefix string, xm int, zoom int, suffix string) *[]byte {
	xmt := []string{"x", "m"}
	img_bytes, ok := FileBytesRead(prefix + xmt[xm] + I2S(zoom) + suffix)
	if ok {
		return img_bytes
	} else {
		za := Constant_ZoomArray()
		scales := Constant_ZoomArray()
		for j := len(za) - 1; j >= 0; j-- {
			scales = append(scales, za[j]/8*3)
		}
		//if xm > 0 {
		//	scales = []int{96, 48, 24, 256, 128, 64}
		//}
		for _, scale := range scales {
			for _, xmv := range xmt {
				fname := prefix + xmv + I2S(scale) + suffix
				if FileExists(fname) {
					img_bytes2, ok2 := FileBytesRead(fname)
					if ok2 {
						return img_bytes2
					}
				}
			}
		}
	}
	return nil
}

func GetExtensionIconName(ftype string, dir bool) string {
	if dir {
		return PREFIX_DRAWONME + FILE_TYPE_FOLDER
	}
	if ftype == "" {
		return PREFIX_DRAWONME + FILE_TYPE_UNKNOWN
	}
	result, ok := TABLE_EXTENSIONS_ICONS_NAMES[ftype]
	if ok {
		return result
	}
	return PREFIX_DRAWONME + FILE_TYPE_UNKNOWN
}

func GetIcon_PixBif_OF(zoom int, iconname string) *gdk.Pixbuf {
	iconset := TABLE_EXTENSIONS_ICONS_SETS[zoom]
	return iconset.BigPixBufs[iconname]
}

func GetIcon_PixBif(zoom int, ftype string, dir bool) *gdk.Pixbuf {
	iconname := GetExtensionIconName(ftype, dir)
	return GetIcon_PixBif_OF(zoom, iconname)
}

func GetIcon_ImageRGBA_OF(zoom int, iconname string) image.Image {
	iconset := TABLE_EXTENSIONS_ICONS_SETS[zoom]
	return iconset.SmallImageRGBA[iconname]
}

func GetIcon_ImageRGBA(zoom int, ftype string, dir bool) image.Image {
	iconname := GetExtensionIconName(ftype, dir)
	return GetIcon_ImageRGBA_OF(zoom, iconname)
}

func GetIcon_ImageFolder(zoom int) *image.RGBA {
	iconset := TABLE_EXTENSIONS_ICONS_SETS[zoom]
	return ImageClone(iconset.FolderImageRGBA)
}

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
