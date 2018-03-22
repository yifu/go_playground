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

type OmittingDirErr struct {
	dirName string
}

func (err OmittingDirErr) Error() string {
	return fmt.Sprintf("Omitting directory %q", err.dirName)
}

type NotADirErr struct {
	paramName string
}

func (err NotADirErr) Error() string {
	return fmt.Sprintf("Target %q is not a directory", err.paramName)
}

func main() {
	flag.Usage = func() {
		fmt.Print("Usage: ", os.Args[0], " sourcefile destdi/\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile1 sourcefile2 ... destdir\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile destfile>\n")
		flag.PrintDefaults()
	}

	paramCount := len(os.Args) - 1

	if paramCount < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// TODO FIXME This checking is not good. We must verify that there is only valid files (this no dir).
	if paramCount == 2 {
		copyFileIntoFile(os.Args[1], os.Args[2])
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

func copyFileIntoFile(srcPath, dstPath string) {
	//fmt.Println("copy into file")
	src, err := os.Open(srcPath)
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer src.Close()

	srcStat, err := src.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_TRUNC, 0 /*perm is useless when O_CREATE is not specified*/)
	//fmt.Println("open trunc")
	if err != nil {
		if os.IsNotExist(err) {
			dst, err = os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, srcStat.Mode().Perm())
			//fmt.Println("open create", dst, ", err", err)
			if err != nil {
				printErr(err)
				os.Exit(2)
			}
			defer dst.Close()
		} else {
			printErr(err)
			os.Exit(1)
		}
	}
	//fmt.Println("after open trunc/create", dst)

	dstStat, err := dst.Stat()
	if err != nil {
		//fmt.Println("dst stat", dst)
		printErr(err)
		os.Exit(2)
	}

	if os.SameFile(srcStat, dstStat) {
		//fmt.Println("same file? yes")
		os.Exit(0)
	}

	//fmt.Println("copy file")
	if _, err := io.Copy(dst, src); err != nil {
		printErr(err)
		os.Exit(2)
	}

	os.Exit(0)
}

func copyFileInDir(srcFileName, destDirName string) {

	fmt.Println("copy file in dir. src=", srcFileName, ", dst=", destDirName)

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
