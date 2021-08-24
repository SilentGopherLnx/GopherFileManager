package pkg_filetools

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"
	//"io/ioutil"
	//"net/url"
)

type DiskPart struct {
	Title        string
	PartName     string
	FSType       string
	Protocol     string
	SpaceTotal   string
	SpaceUsed    string
	SpaceFree    string
	SpacePercent int
	Crypted      bool
	Primary      bool
	MountPath    string
	Model        string
}

func (p *DiskPart) String() string {
	return p.Model + "\t[" + B2S_YN(p.Primary) + "]\t" + p.Title + "\t" + p.PartName + "\t" + p.FSType + "\t" + p.Protocol + "\t" + I2S(p.SpacePercent) + "%\t" + p.SpaceTotal + "\t[" + B2S_YN(p.Crypted) + "]\t" + p.MountPath
}

// df -T -h
// оставить смонтированные и x-gvfs-show
// lsblk -l -n -o "MODEL,PKNAME,KNAME,LABEL,TYPE,FSTYPE,SIZE,RM,HOTPLUG,MOUNTPOINT" //-b bytes   rm- removable  rom - readonly
// mount -v | grep "x-gvfs-show"

//udevadm info --name=/dev/bus/usb/"+usb1+"/"+usb2+" | grep ID_MODEL=   //for MTP

// id -u
// df -T -h "run/user/1000/gvfs/..."

func Linux_DisksGetMounted(local bool, remote bool) []*DiskPart {
	d := []*DiskPart{}

	uid, _, _ := GetPC_UserUidLoginName()

	if local {
		T_ROOT := "ROOT"

		//===== FULL LIST WITH SIZE

		var d_all []*DiskPart
		discs_all, _, _ := ExecCommand("df", "-T", "-h")
		//Prln(discs_all)
		discs_arr := StringSplit(discs_all, "\n")
		for _, disc := range discs_arr[1:] {
			disctrim := disc
			disctrim2 := StringReplace(disctrim, "  ", " ")
			for len(disctrim2) < len(disctrim) {
				disctrim = disctrim2
				disctrim2 = StringReplace(disctrim, "  ", " ")
			}
			disctrim = disctrim2
			t := StringSplit(disctrim, " ")
			//Filesystem             Type      Size  Used Avail Use% Mounted
			//Prln(I2S(len(t)) + "# " + disctrim)
			if len(t) > 6 {
				pr := S2I(StringReplace(t[5], "%", ""))
				ind := StringFind(disc, "%")
				mount := StringPart(disc, ind+2, 0)
				crypt := StringFind(t[0], "_crypt") > 0
				prot := "PART" //"DISK"
				if StringFind(t[1], "tmpfs") > 0 {
					prot = "RAM"
				}
				prim := true
				if t[1] == "fuseblk" {
					prim = false
				}
				d_new := &DiskPart{Title: t[0], PartName: t[0], FSType: t[1], Protocol: prot, SpaceTotal: t[2], SpaceUsed: t[3], SpaceFree: t[4], SpacePercent: pr, Crypted: crypt, Primary: prim, MountPath: mount}
				d_all = append(d_all, d_new)
			}
		}

		//===== ADD REAL DEVICES WITH NAME FIX

		discs_real, _, _ := ExecCommand("lsblk", "-r", "-n", "-b", "-o", "MODEL,PKNAME,KNAME,LABEL,TYPE,FSTYPE,SIZE,HOTPLUG,MOUNTPOINT")
		//Prln(discs_real)
		discs_arr2 := StringSplit(discs_real, "\n")
		tmodel := ""
		for _, disc := range discs_arr2 {
			t := StringSplit(disc, " ")
			/*if len(t) > 0 {
				tmod := t[0]
				if len(tmod) > 0 {
					tmodel = tmod
				}
			}*/
			tmodel = ""
			if len(t) > 8 {
				//if StringDown(t[4]) == "crypt" {
				//tmodel = ""
				for _, disc2 := range discs_arr2 {
					t2 := StringSplit(disc2, " ")
					if len(t2) > 8 {
						if t[1] == t2[2] {
							if len(StringTrim(t2[0])) > 0 {
								tmodel = t2[0]
							} else {
								for _, disc3 := range discs_arr2 {
									t3 := StringSplit(disc3, " ")
									if len(t3) > 8 {
										if t2[1] == t3[2] {
											tmodel = t3[0]
										}
									}
								}
							}
						}
					}
				}
				//}
				mount := t[8]
				if len(mount) > 0 {
					for _, d_real := range d_all {
						if d_real.MountPath == mount {
							if StringDown(t[4]) == "rom" {
								d_real.Protocol = "ROM"
							}
							//if t[7] == "1" {
							d_real.Primary = (t[7] == "0")
							//}
							d_real.Model = tmodel
							d_real.FSType = t[5]
							tt := StringTrim(t[3])
							if len(tt) > 0 {
								d_real.Title = tt
							} else {
								tt = StringTrim(t[2])
								if len(tt) > 0 {
									d_real.Title = tt
								}
							}
							d = append(d, d_real)
						}
					}
				}
			}
		}

		//===== ADD VISIBLE VIRTUAL

		discs_show, _, _ := ExecCommandBash("mount -v | grep x-gvfs-show")
		//Prln(discs_show)
		discs_arr3 := StringSplit(discs_show, "\n")
		for _, disc := range discs_arr3 {
			t := StringSplit(disc, " ")
			for _, d_sh := range d_all {
				if StringFind(t[0], d_sh.PartName) > 0 {
					exist := false
					for _, d_ch := range d {
						if d_ch == d_sh {
							exist = true
						}
					}
					if !exist {
						d = append(d, d_sh)
					}
				}
			}
		}

		// ===== SOME NICE FIX

		for _, ddev := range d {
			dev := ddev.PartName
			devs := StringSplit(dev, "/")
			if len(devs) > 0 {
				dev = devs[len(devs)-1]
				if dev == "" {
					//dev := disc.PartName
				}
			}
			ddev.PartName = StringReplace(dev, "_crypt", "")
			ddev.Model = StringTrim(StringReplace(ddev.Model, "\\x20", " "))
			if ddev.MountPath == "/" {
				ddev.Title = T_ROOT
			}
			if ddev.Title == ddev.PartName {
				tts := StringSplit(ddev.MountPath, "/")
				if len(tts) > 0 {
					tt := StringTrim(tts[len(tts)-1])
					if tt != "" {
						ddev.Title = tt
					}
				}
			}
		}

		//===== NICE SORT

		SortArray(d, func(i, j int) bool {
			if d[i].Primary != d[j].Primary {
				return !CompareBoolLess(d[i].Primary, d[j].Primary)
			}
			/*r1 := d[i].Title == T_ROOT
			r2 := d[j].Title == T_ROOT
			if r1 != r2 {
				return !CompareBoolLess(r1, r2)
			}*/
			if d[i].Protocol != d[j].Protocol {
				return d[i].Protocol < d[j].Protocol
			}
			if d[i].Model != d[j].Model {
				return d[i].Model < d[j].Model
			}
			return false
		})
	}

	if remote {
		/*uid, _, _ := ExecCommand("id", "-u")
		uid = StringTrim(uid)*/
		if len(uid) == 0 {
			uid = "1000"
		}

		remote_dir := "/run/user/" + uid + "/gvfs/"
		files, _ := Folder_ListFiles(remote_dir, false)

		for _, f := range files {
			//if f.IsDir() {
			f_name := f.Name()
			dir_name := remote_dir + f_name + "/"
			alldiscs, _, _ := ExecCommand("df", "-T", "-h", dir_name)

			//Prln(alldiscs)
			discs_arr := StringSplit(alldiscs, "\n")
			for _, disc := range discs_arr[1:2] {
				disctrim := disc
				disctrim2 := StringReplace(disctrim, "  ", " ")
				for len(disctrim2) < len(disctrim) {
					disctrim = disctrim2
					disctrim2 = StringReplace(disctrim, "  ", " ")
				}
				disctrim = disctrim2
				t := StringSplit(disctrim, " ")
				//Filesystem             Type      Size  Used Avail Use% Mounted
				//Prln(I2S(len(t)) + "# " + disctrim)
				if len(t) > 6 {
					title := f_name
					pr := S2I(StringReplace(t[5], "%", ""))
					//ind := StringFind(disc, "%")
					//mount := StringPart(disc, ind+2, 0) + "/" + f_name
					mount := dir_name
					dev := t[0]
					prot := "Other"
					indp := StringFind(f_name, ":")
					mdl := ""
					crypt := false
					if indp > 0 {
						title = StringPart(f_name, indp+1, 0)
						//q, qerr := UrlQueryParse(StringReplace(StringReplace(title, ";", "#"), ",", ";"))
						q, qerr := UrlQueryParse(StringReplace(title, ",", "&"))
						prot = StringDown(StringPart(f_name, 0, indp-1))
						if StringFind(prot, "smb") == 1 {
							prot = "SMB"
							if qerr == nil {
								mdl = StringTrim(q.Get("server"))
								share := StringTrim(q.Get("share"))
								if len(share) > 0 {
									title = "[" + share + "]"
								}
							} else {
								//Prln("title=" + title)
							}
						}
						if StringFind(prot, "dav") == 1 {
							prot = "WEBDAV"
							if StringFind(title, "ssl=true") > 0 {
								crypt = true
							}
							if qerr == nil {
								mdl = StringTrim(q.Get("host"))
								if len(mdl) > 0 {
									user := StringTrim(q.Get("user"))
									if len(user) > 0 {
										title = "[" + StringUp(user) + "]"
									}
								}
							}
						}
						if StringFind(prot, "ftp") == 1 {
							prot = "FTP"
							if StringFind(title, "ssl=true") > 0 {
								crypt = true
							}
							if qerr == nil {
								mdl = StringTrim(q.Get("host"))
								if len(mdl) > 0 {
									user := StringTrim(q.Get("user"))
									if len(user) > 0 {
										title = "[" + StringUp(user) + "]"
									} else {
										title = "[" + StringUp(mdl) + "]"
									}
								}
							}
							Prln("+[[[[" + mount + "]]]]" + f_name)
						}
						if StringFind(prot, "ssh") == 1 { // NOT TESTED
							prot = "SSH"
							crypt = true
							if qerr == nil {
								mdl = StringTrim(q.Get("host"))
								if len(mdl) == 0 {
									mdl = StringTrim(q.Get("server"))
								}
								if len(mdl) > 0 {
									user := StringTrim(q.Get("user"))
									if len(user) > 0 {
										title = "[" + StringUp(user) + "]"
									}
								}
							}
						}
						//===========
						if StringFind(prot, "mtp") == 1 || StringFind(prot, "gphoto") > 0 { // "gphoto2"
							if StringFind(prot, "mtp") == 1 {
								prot = "MTP"
							} else if StringFind(prot, "gphoto") == 1 {
								prot = "PTP"
							} else {
								prot = StringUp(prot)
							}
							indh := StringFind(title, "=")
							if len(title) > indh {
								title = StringPart(title, indh+1, 0)
							}
							//host=%5Busb%3A001%2C018%5D
							//mtp://[usb:001,006]/
							title = UrlQueryUnescape(title)
							ind1 := StringFind(title, ":") //%3A")
							ind2 := StringFind(title, ",") //"%2C")
							ind3 := StringFind(title, "]") //"%5D")
							if ind1 > 1 && ind2 > ind1 && ind3 > ind2 {
								usb1 := StringPart(title, ind1+1, ind2-1)
								usb2 := StringPart(title, ind2+1, ind3-1)
								dev_name, _, _ := ExecCommandBash("udevadm info --name=/dev/bus/usb/" + usb1 + "/" + usb2 + " | grep ID_MODEL=")
								dev_name = StringTrim(dev_name)
								if len(dev_name) > 0 {
									mdl = dev_name
									ind0 := StringFind(mdl, "=")
									if ind0 > 1 {
										mdl = StringPart(mdl, ind0+1, 0)
										title = "{usb=" + usb1 + "," + usb2 + "}"
									} else {
										title = "{" + title + "}"
									}
								}
							} else {
								title = "{" + title + "}"
							}
						}
					}
					d_new := &DiskPart{Title: title, PartName: dev, FSType: t[1], Protocol: prot, SpaceTotal: t[2], SpaceUsed: t[3], SpaceFree: t[4], SpacePercent: pr, Crypted: crypt, Primary: false, MountPath: mount, Model: mdl}
					d = append(d, d_new)
				}
			}
			//}
		}
	}

	return d
}

func Linux_DisksGetAllLocal() []*DiskPart {
	//lsblk  --noheadings --raw | awk '{print substr($0,0,4)}' | uniq -c | grep 1 | awk '{print "/dev/"$2}'
	//cat /proc/partitions
	//sfdisk -l

	//lsblk -l -n -o "MODEL,PKNAME,KNAME,LABEL,TYPE,FSTYPE,SIZE,RM,HOTPLUG,MOUNTPOINT" - best
	d := []*DiskPart{}
	d = d
	//return d
	return Linux_DisksGetMounted(true, false)
}
