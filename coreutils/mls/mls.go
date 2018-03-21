package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
)

var showHidden, showInode, longListingFmt = false, false, false

func main() {
	flag.BoolVar(&showHidden, "a", false, "Show hidden files.")
	flag.BoolVar(&showInode, "i", false, "Show inode numbers.")
	flag.BoolVar(&longListingFmt, "l", false, "Enable long listing format.")

	flag.Parse()

	var param string
	if len(flag.Args()) == 0 {
		param = "."
	} else {
		param = flag.Arg(0)
	}

	fileInfo, err := os.Stat(param)
	if err != nil {
		fmt.Println(os.Args[0]+":", err.Error())
		os.Exit(1)
	}

	switch {
	case fileInfo.IsDir():
		processDir(param)
	default:
		processFile(fileInfo)
	}
}

func processDir(dirPath string) {
	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	for _, fileInfo := range infos {
		if strings.HasPrefix(fileInfo.Name(), ".") && !showHidden {
			continue
		}

		processFile(fileInfo)
	}
}

func processFile(fileInfo os.FileInfo) {
	if showInode {
		fmt.Print(uint64(fileInfo.Sys().(*syscall.Stat_t).Ino), " ")
	}

	if longListingFmt {
		stat := fileInfo.Sys().(*syscall.Stat_t)
		dateMtime := time.Unix(stat.Mtim.Unix()).Format("2 Jan")
		hourMtime := time.Unix(stat.Mtim.Unix()).Format("15:04")
		fmt.Printf("%v %4d %v %v %5v %6v %v ",
			fileInfo.Mode(), stat.Nlink, stat.Uid, stat.Gid,
			stat.Size, dateMtime, hourMtime)
	}

	fmt.Println(fileInfo.Name())
}
