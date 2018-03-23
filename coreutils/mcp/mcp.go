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
// 7a- mcp -r /a /b when b does not already exists
// 7b- mcp -r /a /b when b already exists
// Every time: check the resulting mode for every new file/dir.

// TODO Next steps: implement -r option.

// TODO Replace those structs with fmt.Errorf(fmt, "")

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

type NoSuchFileOrDirErr struct {
	paramName string
}

func (err NoSuchFileOrDirErr) Error() string {
	return fmt.Sprintf("%q No such file or directory", err.paramName)
}

func main() {
	flag.Usage = func() {
		fmt.Print("Usage: ", os.Args[0], " sourcefile destdir/\n")
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
			if os.IsNotExist(err) {
				printErr(NoSuchFileOrDirErr{paramName: srcName})
				os.Exit(1)
			} else {
				printErr(err)
				os.Exit(2)
			}
		}
		if srcInfo.IsDir() {
			printErr(OmittingDirErr{dirName: srcName})
			os.Exit(1)
		}
		// TODO When the second parameter is a dir, it must work.
		copyFileIntoFile(srcName, dstName)
		os.Exit(0)
	}

	destDir := os.Args[len(os.Args)-1]
	if !isDir(destDir) {
		printErr(NotADirErr{paramName: destDir})
		os.Exit(1)
	}

	copyFiles(destDir, os.Args[1:len(os.Args)-1]...)
}

func copyFiles(destDir string, srcList ...string) {
	for i, param := range srcList {
		fileInfo, err := os.Stat(param)
		if err != nil {
			if os.IsNotExist(err) {
				// TODO We must exit with error number 1 while still processing the rest of the params.
				printErr(NoSuchFileOrDirErr{paramName: param})
				continue
			} else {
				printErr(err)
				os.Exit(2)
			}
		}

		if fileInfo.IsDir() {
			printErr(OmittingDirErr{dirName: param})
			continue
		}

		_, filename := filepath.Split(param)
		if findFilename(srcList[:i], filename) {
			// TODO We must exit with status code 1 in the end, but still process the rest of the params.
			fmt.Printf("%v will not overwrite '' with %q\n", os.Args[0], filename)
			continue
		}

		copyFileIntoDir(param, destDir)
	}

}

func isDir(destDir string) bool {
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

	return dstStat.IsDir()
}

func findFilename(filepaths []string, filename string) bool {
	for _, param := range filepaths {
		_, name := filepath.Split(param)
		if filename == name {
			return true
		}
	}

	return false
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
