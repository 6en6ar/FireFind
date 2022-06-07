package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

var dir = "DecompiledApk"

func DecompileWithApk(apkFile string) string {
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	u := uuid.New()
	uuid := strings.Replace(u.String(), "-", "", -1)
	apkDir := dir + "/" + uuid
	// execute apktool
	out, err := exec.Command("apktool", "d", apkFile, "--output", apkDir).Output()
	if err != nil {
		color.Red("[ - ] APKTool failed with " + err.Error())
		os.Exit(1)

	}
	color.Green("[ + ] Decompiling APK file ...")
	color.Yellow(string(out))
	color.Green("[ + ] Successfully decompiled the app using ApkTool...")
	return apkDir
}
func FindFirebaseInstance(apkDir string) {
	FirebaseRe, _ := regexp.Compile("https*(.*?).firebaseio.com")
	err := filepath.Walk(apkDir, func(filename string, info os.FileInfo, err error) error {
		fhandle, err := os.Open(filename)
		//fmt.Println("[ + ] Checking file -->" + filename)
		f := bufio.NewReader(fhandle)

		if err != nil {
			color.Red("[ - ] Error opening file " + filename)
		}
		defer fhandle.Close()

		buf := make([]byte, 1024)
		for {
			buf, _, err = f.ReadLine()
			if err != nil {
				break
			}

			s := string(buf)
			if FirebaseRe.MatchString(s) {
				url := FirebaseRe.FindStringSubmatch(s)
				Firebaseinstances = append(Firebaseinstances, url[0])

			}
		}
		return nil
	})
	if err != nil {
		color.Red("[ - ] Error traversing directory ..." + err.Error())
	}
}

func CheckInstance(instance string) {
	res, err := http.Get(instance + "/.json")
	if err != nil {
		log.Fatalln(err)
	}
	if res.StatusCode == 404 {
		color.Red("[ - ] Instance " + instance + " doesn't exist!")
		os.Exit(1)
	}
	if res.StatusCode == 403 {
		color.Green("[ + ] Instance " + instance + " seems to be secured!")
		os.Exit(1)
	}
	if res.StatusCode == 423 {
		color.Red("[ - ] Instance " + instance + " seems to be deleted!")
		os.Exit(1)
	}
	if res.StatusCode == 200 {
		color.Green("[ + ] Found Misconfigured Firebase instance --> " + instance)
	}
}
func usage() {
	color.Blue("usage : -a=<PATH-TO-APK-FILE>\n")
	os.Exit(0)
}

var banner = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣆⢳⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣾⣷⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠠⣄⠀⢠⣿⣿⣿⣿⡎⢻⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢸⣧⢸⣿⣿⣿⣿⡇⠀⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣾⣿⣿⣿⣿⠃⠀⢸⣿⣿⣿⣿⣿⣿⠀⣄⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⣿⣿⠏⠀⠀⣸⣿⣿⣿⣿⣿⡿⢀⣿⡆⠀
⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣿⣿⠃⠀⠀⠀⣿⣿⣿⣿⣿⣿⠇⣼⣿⣿⡄
⠀⢰⠀⠀⣴⣿⣿⣿⣿⣿⣿⡿⠁⠀⠀⠀⢠⣿⣿⣿⣿⣿⡟⣼⣿⣿⣿⣧
⠀⣿⡀⢸⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⣸⡿⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿
⠀⣿⣷⣼⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⢹⠃⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿
⡄⢻⣿⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⠇
⢳⣌⢿⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠻⣿⣿⣿⣿⣿⠏⠀
⠀⢿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⠋⣠⠀
⠀⠈⢻⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣵⣿⠃⠀
⠀⠀⠀⠙⢿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⡿⠃⠀⠀
⠀⠀⠀⠀⠀⠙⢿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⡿⠋⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠛⠿⣿⣦⣀⠀⠀⠀⠀⢀⣴⠿⠛⠁⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠓⠂⠀⠈⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
FireFind ... Coded by 6en6ar >:)
`
var Firebaseinstances []string

func main() {
	color.Red(banner)
	apk := flag.String("a", "", "Path to an APK file you wish to test")
	flag.Usage = usage
	flag.Parse()
	if flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}
	apkDir := DecompileWithApk(*apk)
	FindFirebaseInstance(apkDir)
	for i := 0; i < len(Firebaseinstances); i++ {
		color.Green("[ + ] Found instance --> " + Firebaseinstances[i])
		color.Green("[ + ] Checkin accessibility ...")
		CheckInstance(Firebaseinstances[i])
	}

}
