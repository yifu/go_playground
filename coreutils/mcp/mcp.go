package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TODO Implement a test suite of some sort.
// TODO implémenter cp -r
// TODO implémenter cp -f

func main() {
	setUsage()
	checkArgsCount()

	dst, srcs := parseCmdLine()
	if isDstDir := isDir(dst); len(srcs) > 1 && !isDstDir {
		printErr(NotADirErr{dst})
		os.Exit(1)
	} else {
		var errs []error
		srcs, errs = filterSrcList(dst, srcs)

		for _, src := range srcs {
			if isDstDir {
				_, fileName := filepath.Split(src)
				target := filepath.Join(dst, fileName)
				cp(src, target)
			} else {
				cp(src, dst)
			}
		}

		for _, err := range errs {
			fmt.Println(err.Error())
		}

		if len(errs) > 0 {
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func parseCmdLine() (string, pathList) {
	dst := os.Args[len(os.Args)-1]
	srcs := os.Args[1 : len(os.Args)-1]
	return dst, srcs
}

func isDir(f string) bool {
	if fi, err := os.Stat(f); err != nil {
		return false
	} else {
		return fi.IsDir()
	}
}

func printErr(e error) {
	fmt.Print(os.Args[0], ": ", e.Error(), "\n")
}

func cp(src, dst string) {
	if sameFile(src, dst) {
		return
	}

	srcf, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer srcf.Close()

	srcfi, err := srcf.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
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

func filterSrcList(dst string, srcs pathList) (oks pathList, errors []error) {
	for _, param := range srcs {
		fileInfo, err := os.Stat(param)
		if err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, NoSuchFileOrDirErr{paramName: param})
				continue
			} else {
				printErr(err)
				os.Exit(2)
			}
		}

		if fileInfo.IsDir() {
			errors = append(errors, OmittingDirErr{dirName: param})
			continue
		}

		_, fileName := filepath.Split(param)
		if oks.contains(fileName) {
			errors = append(errors, WillNotOverwriteErr{paramName: param, alreadyCopied: filepath.Join(dst, fileName)})
			continue
		}

		oks = append(oks, param)
	}
	return
}

type pathList []string

func (paths pathList) contains(path string) bool {
	_, fn := filepath.Split(path)
	for _, p := range paths {
		_, n := filepath.Split(p)
		if fn == n {
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
