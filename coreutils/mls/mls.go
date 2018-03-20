package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
)

func main() {
	var showHidden = false
	flag.BoolVar(&showHidden, "a", false, "Show hidden files.")

	var showInode = false
	flag.BoolVar(&showInode, "i", false, "Show inode numbers.")

	flag.Parse()

	infos, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	for _, fileinfo := range infos {
		if strings.HasPrefix(fileinfo.Name(), ".") && !showHidden {
			continue
		}

		if showInode {
			fmt.Print(uint64(fileinfo.Sys().(*syscall.Stat_t).Ino), " ")
		}

		fmt.Println(fileinfo.Name())

	}
}
