package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	for _, fileinfo := range infos {
		fmt.Println(fileinfo.Name())
	}
}
