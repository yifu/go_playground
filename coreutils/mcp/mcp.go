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

	if paramCount == 2 {
		srcName, dstName := os.Args[1], os.Args[2]
		srcInfo, err := os.Stat(srcName)
		if err != nil {
			printErr(err)
			os.Exit(2)
		}
		if srcInfo.IsDir() {
			printErr(OmittingDirErr{dirName: srcName})
			os.Exit(1)
		}
		copyFileIntoFile(srcName, dstName)
		os.Exit(0)
	}

	destDir := os.Args[len(os.Args)-1]
	dstInfo, err := os.Open(destDir)
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	dstStat, err := dstInfo.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	if !dstStat.IsDir() {
		printErr(NotADirErr{paramName: destDir})
		os.Exit(1)
	}

	for i, filename := range os.Args {
		if i == 0 || i == len(os.Args)-1 {
			continue
		}

		// TODO We must check if the filename has already been copied into the dest dir during this mcp execution. When it's been the case, we skip after priting a message.
		copyFileIntoDir(filename, destDir)
	}
}

func printErr(e error) {
	fmt.Print(os.Args[0], ": ", e.Error(), "\n")
}

func copyFileIntoFile(srcPath, dstPath string) {
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

	var dst *os.File
	dstStat, err := os.Stat(dstPath)
	if err != nil {
		if os.IsNotExist(err) {
			dst, err = os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, srcStat.Mode().Perm())
			if err != nil {
				printErr(err)
				os.Exit(2)
			}
			defer dst.Close()
		} else {
			printErr(err)
			os.Exit(2)
		}
	} else {
		if os.SameFile(srcStat, dstStat) {
			os.Exit(0)
		}

		dst, err = os.OpenFile(dstPath, os.O_WRONLY|os.O_TRUNC, 0 /*perm is useless when O_CREATE is not specified*/)
		if err != nil {
			printErr(err)
			os.Exit(2)
		}
	}

	if _, err := io.Copy(dst, src); err != nil {
		printErr(err)
		os.Exit(2)
	}
}

func copyFileIntoDir(srcPath, destDirPath string) {
	_, srcFileName := filepath.Split(srcPath)
	dstFilePath := filepath.Join(destDirPath, srcFileName)

	copyFileIntoFile(srcPath, dstFilePath)
}