package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	var showHidden = false
	flag.BoolVar(&showHidden, "a", false, "Show hidden files.")
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
		fmt.Println(fileinfo.Name())
	}
}
