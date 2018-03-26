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

func openSrc(src string) (*os.File, os.FileInfo, error) {
	srcf, err := os.OpenFile(src, os.O_RDONLY, 0)
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
	return srcf, srcfi, nil
}

func mkDst(dst, src string) string {
	// Open dst, readable only. We will not read into it. It's just to check it exists.
	if dstf, err := os.OpenFile(dst, os.O_RDONLY, 0); err != nil {
		if !os.IsNotExist(err) {
			printErr(err)
			os.Exit(2)
		}
		// Nothing to do: dst is a filename,
		// which does not exist yet.
		// The user is asking to copy and rename at the same time.
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

		// Open src
		srcf, srcfi, err := openSrc(src)
		if err != nil {
			printErr(err)
			exitCode = 1
			continue
		}
		if srcfi.IsDir() {
			printErr(OmittingDirErr{srcfi.Name()})
			exitCode = 1
			continue
		}
		if oks.contains(dst) {
			printErr(WillNotOverwriteErr{src, dst})
			exitCode = 1
			continue
		}
		if dstfi, err := os.Stat(dst); err != nil {
			if os.IsNotExist(err) {
			} else {
				printErr(err)
				os.Exit(2)
			}
		} else if os.SameFile(dstfi, srcfi) {
			// Nothing to copy.
			continue
		}
		// Open dst, writable this time
		dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
		if err != nil {
			printErr(err)
			os.Exit(2)
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
