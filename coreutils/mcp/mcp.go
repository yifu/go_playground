package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	flag.Usage = func() {
		fmt.Print("Usage: ", os.Args[0], " sourcefile destdi/\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile1 sourcefile2 ... destdir\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile destfile>\n")
		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	fileInfo, err := os.Stat(os.Args[2])
	if err != nil {
		fmt.Print(os.Args[0], ": ", err.Error(), "\n")
		os.Exit(2)
	}

	if !fileInfo.IsDir() {
		if len(os.Args) == 3 {

			// TODO Check both param points to the same file using os.SameFile().
			if os.Args[1] == os.Args[2] {
				os.Exit(0)
			}

			// TODO Copy the src file into the dst file open(chemin, O_WRONLY | O_TRUNC) ou open(chemin,  O_WRONLY | O_CREAT,  mode)
		}
		flag.Usage()
		os.Exit(1)
	}

	destDir := os.Args[len(os.Args)-1]

	for i, filename := range os.Args {
		if i == 0 || i == len(os.Args)-1 {
			continue
		}

		copyFileInDir(filename, destDir)
	}
}

func printErr(e error) {
	fmt.Print(os.Args[0], ": ", e.Error(), "\n")
}

func copyFileInDir(srcFileName, destDirName string) {
	src, err := os.Open(srcFileName)
	if err != nil {
		fmt.Print(os.Args[0], ": ", err.Error())
		os.Exit(2)
	}
	defer src.Close()

	srcFileInfo, err := src.Stat()
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	perm := srcFileInfo.Mode().Perm()

	_, filename := filepath.Split(filepath.Clean(srcFileName))

	dstFileName := filepath.Clean(destDirName) + string(filepath.Separator) + filename

	dst, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	defer dst.Close()
	fmt.Println("copy..", src, ", ", dst)

	if _, err := io.Copy(dst, src); err != nil {
		printErr(err)
		os.Exit(2)
	}
}
