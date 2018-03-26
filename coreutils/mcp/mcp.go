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

type pathList []string

// TODO Placer des defer f.Close() partout où nécessaire.

func main() {
	setUsage()
	checkArgsCount()
	exitCode := 0

	dst, srcs := parseCmdLine()
	oks := make(pathList, 0, len(srcs))

	var dir string
	if len(srcs) > 1 {
		// When len(srcs) > 1, then dst must be a dir. We check that it is the case:
		if !checkIsDir(dst) {
			printErr(NotADirErr{dst})
			os.Exit(1)
		}
		dir = dst
	}

	for _, src := range srcs {
		switch {
		case len(srcs) == 1:
			dst = mkDst(dst, src)
		case len(srcs) > 1:
			_, fn := filepath.Split(src)
			dst = filepath.Join(dir, fn)
		}
		dstf, srcf, err := openFiles(dst, src, oks)
		if err != nil {
			if _, ok := err.(SameFileErr); ok {
				// When src and dst are the same files, there is no copy goind on.
				// We just keep going on with the next src.
			} else {
				printErr(err)
				exitCode = 1
			}
			continue
		}
		if _, err := io.Copy(dstf, srcf); err != nil {
			printErr(err)
			os.Exit(2)
		}
		oks = append(oks, dst)
	}
	os.Exit(exitCode)
}

func parseCmdLine() (string, pathList) {
	dst := os.Args[len(os.Args)-1]
	srcs := os.Args[1 : len(os.Args)-1]
	return dst, srcs
}

func printErr(e error) {
	fmt.Print(os.Args[0], ": ", e.Error(), "\n")
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

func checkIsDir(dst string) bool {
	dirfi, err := os.Stat(dst)
	return err == nil && dirfi.IsDir()
}

func mkDst(dst, src string) string {
	// Open dst, readable only. We will not read into it. It's just to check it exists.
	// TODO Mhh.. We may actually stat() the file directly. Why open it? Not sure now.
	if dstf, err := os.OpenFile(dst, os.O_RDONLY, 0); err != nil {
		if os.IsNotExist(err) {
			// Nothing to do: dst is a filename,
			// which does not exist yet.
			// The user is asking to copy and rename at the same time.
		} else {
			printErr(err)
			os.Exit(2)
		}
	} else {
		dstfi, err := dstf.Stat()
		if err != nil {
			printErr(err)
			os.Exit(2)
		}
		if dstfi.IsDir() {
			return filepath.Join(dst, src)
		}
		// Nothing to do: dst is a filename, we must copy into it directly.
	}
	return dst
}

func (paths pathList) contains(path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func openFiles(dst, src string, oks pathList) (dstf, srcf *os.File, err error) {
	// Open src
	srcf, err = os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, NoSuchFileOrDirErr{src}
		}
		printErr(err)
		os.Exit(2)
	}
	srcfi, err := srcf.Stat()
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	// Various checks on src.
	if srcfi.IsDir() {
		return nil, nil, OmittingDirErr{src}
	}
	if oks.contains(dst) {
		return nil, nil, WillNotOverwriteErr{src, dst}
	}
	// Stat(dst) first, then only opening it, to avoid creating an empty file.
	if dstfi, err := os.Stat(dst); err != nil {
		if !os.IsNotExist(err) {
			printErr(err)
			os.Exit(2)
		}
	} else if os.SameFile(dstfi, srcfi) {
		// Nothing to copy.
		return nil, nil, SameFileErr{dst, src}
	}
	dstf, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
	if err != nil {
		printErr(err)
		os.Exit(2)
	}
	return
}

// TODO Replace those structs with fmt.Errorf(fmt, "")

type SameFileErr struct {
	dst, src string
}

func (err SameFileErr) Error() string {
	return fmt.Sprintf("Same Files %q %q", err.dst, err.src)
}

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
