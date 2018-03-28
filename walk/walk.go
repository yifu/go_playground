package main

import (
	"fmt"
	"path/filepath"
	"os"
)

func main() {
	fmt.Println("hello walk")
	filepath.Walk("toto///titi//", walk)
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	fmt.Println("path:", path)
	fmt.Println("dir:", filepath.Dir("tata"))
	return nil
}