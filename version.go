package main

const app_version_manager = "0.4.2" // for automatc update check

func AppVersion() string {
	return app_version_manager
}

func AppAuthor() string {
	return "SilentGopherLnx (2019-2021)"
}

func AppMail() string {
	return "silentgopherlnx@gmail.com"
}

func AppRepository() string {
	return "github.com/SilentGopherLnx/GopherFileManager"
}

func AppAboutMore() string {
	return "Check you have installed:" +
		"\n" +
		"github.com/SilentGopherLnx/FileMoverGui" +
		"\n" +
		"sudo apt-get install ffmpeg" +
		"\n" +
		"sudo apt-get install smbclient"
}

func UrlLastVerison_Manager() string {
	return "https://raw.githubusercontent.com/SilentGopherLnx/GopherFileManager/master/version.go"
}

func UrlLastVerison_Mover() string {
	return "https://raw.githubusercontent.com/SilentGopherLnx/FileMoverGUI/master/version.go"
}
