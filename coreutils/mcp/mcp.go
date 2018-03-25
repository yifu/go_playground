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
// TODO Implémenter un type spécial: path. Faire remonter ces définition en haut du fichier.

func main() {
	setUsage()
	checkArgsCount()
	exitCode := 0

	dst, srcs := parseCmdLine()
	if isDstDir := isDir(dst); len(srcs) > 1 && !isDstDir {
		printErr(NotADirErr{dst})
		exitCode = 1
	} else {
		oks := make(pathList, 0, len(srcs))

		for _, src := range srcs {
			target := mkTargetPath(dst, src, isDstDir)

			if oks.contains(target) {
				printErr(WillNotOverwriteErr{src: src, dst: target})
				exitCode = 1
			} else if err := cp(src, target); err != nil {
				printErr(err)
				exitCode = 1
			} else {
				oks = append(oks, target)
			}
		}
	}
	os.Exit(exitCode)
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

func mkTargetPath(dst, src string, isDstDir bool) string {
	if isDstDir {
		_, fn := filepath.Split(src)
		return filepath.Join(dst, fn)
	} else {
		return dst
	}
}

func printErr(e error) {
	fmt.Print(os.Args[0], ": ", e.Error(), "\n")
}

func cp(src, dst string) error {

	srcf, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return NoSuchFileOrDirErr{src}
		} else {
			printErr(err)
			os.Exit(2)
		}
	}
	defer srcf.Close()

	srcfi, err := srcf.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	if srcfi.IsDir() {
		return OmittingDirErr{src}
	}

	srcperm := srcfi.Mode().Perm()
	dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcperm)
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	defer dstf.Close()

	dstfi, err := dstf.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}

	if os.SameFile(dstfi, srcfi) {
		// Nothing to perform. One side-effect: atime is modified. It is acceptable?
		return nil
	}

	if _, err := io.Copy(dstf, srcf); err != nil {
		printErr(err)
		os.Exit(2)
	}
	return nil
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

type pathList []string

func (paths pathList) contains(path string) bool {
	for _, p := range paths {
		if p == path {
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
	src, dst string
}

func (err WillNotOverwriteErr) Error() string {
	return fmt.Sprintf("Will not overwrite %q with %q", err.dst, err.src)
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
