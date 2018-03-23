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
		fmt.Print("Usage: ", os.Args[0], " sourcefile destdir/\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile1 sourcefile2 ... destdir\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile destfile>\n")
		flag.PrintDefaults()
	}

	checkArgsCount()

	target := os.Args[len(os.Args)-1]
	paramList := os.Args[1 : len(os.Args)-1]

	if isNotExist(target) {
		exitCode := processNonExistingTarget(target, paramList)
		os.Exit(exitCode)
	}

	dstDir := target
	if !isDir(dstDir) {
		printErr(NotADirErr{dstDir})
		os.Exit(1)
	}

	errorList := copyFiles(dstDir, paramList...)

	for _, err := range errorList {
		fmt.Println(err.Error())
	}

	if len(errorList) > 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
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

func checkArgsCount() {
	paramCount := len(os.Args[1:])
	if paramCount < 2 {
		flag.Usage()
		os.Exit(1)
	}
}

func isNotExist(file string) bool {
	_, err := os.Stat(file)
	return os.IsNotExist(err)
}

func processNonExistingTarget(target string, paramList []string) int {
	if len(paramList) == 1 {
		srcFileName, destFileName := paramList[0], target
		copyFileIntoFile(srcFileName, destFileName)
		return 0
	} else {
		printErr(NotADirErr{paramName: target})
		return 1
	}
}

func copyFiles(dstDir string, srcList ...string) []error {
	errorList := make([]error, 0)
	for i, param := range srcList {
		fileInfo, err := os.Stat(param)
		if err != nil {
			if !os.IsNotExist(err) {
				printErr(err)
				os.Exit(2)
			}
			errorList = append(errorList, NoSuchFileOrDirErr{paramName: param})
			continue
		}

		if fileInfo.IsDir() {
			errorList = append(errorList, OmittingDirErr{dirName: param})
			continue
		}

		_, fileName := filepath.Split(param)
		if findFilename(srcList[:i], fileName) {
			errorList = append(errorList, WillNotOverwriteErr{paramName: param, alreadyCopied: filepath.Join(dstDir, fileName)})
			continue
		}

		target := filepath.Join(dstDir, fileName)
		copyFileIntoFile(param, target)
	}
	return errorList
}

func findFilename(filepaths []string, filename string) bool {
	_, filename = filepath.Split(filename)
	for _, param := range filepaths {
		_, name := filepath.Split(param)
		if filename == name {
			return true
		}
	}

	return false
}

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

type WillNotOverwriteErr struct {
	paramName, alreadyCopied string
}

func (err WillNotOverwriteErr) Error() string {
	return fmt.Sprintf("Will not overwrite %q with %q\n", err.alreadyCopied, err.paramName)
}

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
