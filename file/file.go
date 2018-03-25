package main

import (
	"fmt"
	"os"
)

func main() {
	if _, err := os.OpenFile("./file.go/toto.txt", os.O_RDONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
	} else {
		fmt.Println("ok")
	}
}