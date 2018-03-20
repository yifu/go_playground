package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	for _, fileinfo := range infos {
		if strings.HasPrefix(fileinfo.Name(), ".") {
			continue
		}
		fmt.Println(fileinfo.Name())
	}
}
