package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	setUsage()
	checkArgsCount()

	dst, srcs := cmdLine()
	if len(srcs) == 1 {
		copyFileIntoFile(srcs[0], dst)
	} else {
		processCopyingMultipleFiles(dst, srcs)
	}
	os.Exit(0)
}

func cmdLine() (string, []string) {
	dst := os.Args[len(os.Args)-1]
	srcs := os.Args[1 : len(os.Args)-1]
	return dst, srcs
}

func processCopyingMultipleFiles(dst string, srcs []string) {
	if !isDir(dst) {
		printErr(NotADirErr{dst})
		os.Exit(1)
	}

	var errs []error
	srcs, errs = filterSrcList(dst, srcs)

	for _, src := range srcs {
		_, fileName := filepath.Split(src)
		target := filepath.Join(dst, fileName)
		copyFileIntoFile(src, target)
	}

	for _, err := range errs {
		fmt.Println(err.Error())
	}

	if len(errs) > 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func isDir(dst string) bool {
	dstInfo, err := os.Open(dst)
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

func openSrc(src string) (srcf *os.File, srcfi os.FileInfo) {
	var err error
	srcf, err = os.Open(src)
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	srcfi, err = srcf.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	return
}

func copyFileIntoFile(srcPath, dstPath string) {
	if sameFile(srcPath, dstPath) {
		return
	}

	src, err := os.OpenFile(srcPath, os.O_RDONLY, 0)
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

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcStat.Mode().Perm())
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		printErr(err)
		os.Exit(2)
	}
}

func setUsage() {
	flag.Usage = func() {
		fmt.Print("Usage: ", os.Args[0], " sourcefile destdir/\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile1 sourcefile2 ... destdir\n")
		fmt.Print("Usage: ", os.Args[0], " sourcefile destfile>\n")
		flag.PrintDefaults()
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

func sameFile(a, b string) bool {
	afi, err := os.Stat(a)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		printErr(err)
		os.Exit(2)
	}

	bfi, err := os.Stat(b)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		printErr(err)
		os.Exit(2)
	}

	return os.SameFile(afi, bfi)
}

func filterSrcList(dstDir string, srcList []string) (oks []string, errors []error) {
	for i, param := range srcList {
		fileInfo, err := os.Stat(param)
		if err != nil {
			if !os.IsNotExist(err) {
				printErr(err)
				os.Exit(2)
			}
			errors = append(errors, NoSuchFileOrDirErr{paramName: param})
			continue
		}

		if fileInfo.IsDir() {
			errors = append(errors, OmittingDirErr{dirName: param})
			continue
		}

		_, fileName := filepath.Split(param)
		if findFilename(srcList[:i], fileName) {
			errors = append(errors, WillNotOverwriteErr{paramName: param, alreadyCopied: filepath.Join(dstDir, fileName)})
			continue
		}

		oks = append(oks, param)
	}
	return
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
	return fmt.Sprintf("Will not overwrite %q with %q", err.alreadyCopied, err.paramName)
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
