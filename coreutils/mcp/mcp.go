package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TESTS:
// 1a- mcp file1 file2 (file2 does not exist)
// 1b- mcp file1 file2 (file2 already exist)
// 1c- mcp file1 file1 (does nothing, check with mtime)
// 1d- mcp file1 file2 *then* mcp -f file1 file2 (with file2 being beforehand chmod to be not openable. So mcp must unlink() file2)
// 2- mcp file1 ./file2 somewhere/file3 ../file4 ..///../file5 /tmp/file6 dir
// 3- mcp file1 file2 ./dir
// 4- mcp file2 file2 ../dir
// 5- mcp file1 file2 /tmp/dir
// 6- mcp file1 file2 .//.././../dir
// 7- mcp -r /a /b when b already exists
// Every time: check the resulting mode for every new file/dir.

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

	// TODO FIXME This checking is not good. We must verify that there is only valid files (this no dir).
	if allParametersAreFiles() {
		processAllParamAreFiles()
		os.Exit(1)
	}

	destDir := os.Args[len(os.Args)-1]

	// TODO Check this dir does exist. If it does not exist, then we must abort. Special stuff are done only when there are only two parameters.

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

func allParametersAreFiles() bool {
	allRegular := true

	for i, param := range os.Args {
		if i == 0 {
			continue
		}

		fileInfo, err := os.Stat(param)
		if err != nil {
			printErr(err)
			os.Exit(2)
		}

		if !fileInfo.Mode().IsRegular() {
			fmt.Println("not regular")
			allRegular = false
		} else {
			fmt.Println("regular")
		}
	}

	return allRegular
}

func processAllParamAreFiles() {
	if len(os.Args) > 3 {
		flag.Usage()
		os.Exit(1)
	}

	fileInfo1, err := os.Stat(os.Args[1])
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	fileInfo2, err := os.Stat(os.Args[2])
	if err != nil {
		printErr(err)
		// TODO
		os.Exit(1)
	}

	if os.SameFile(fileInfo1, fileInfo2) {
		os.Exit(0)
	}

	// TODO Copy the src file into the dst file open(chemin, O_WRONLY | O_TRUNC) ou open(chemin,  O_WRONLY | O_CREAT,  mode)
	src, err := os.Open(os.Args[1])
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer src.Close()

	dst, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_TRUNC, fileInfo1.Mode().Perm())
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer dst.Close()

	io.Copy(dst, src)
	os.Exit(0)
}

func copyFileInDir(srcFileName, destDirName string) {

	fmt.Println("src=", srcFileName, ", dst=", destDirName)

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

	fmt.Println("dstfilename = ", dstFileName)
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
