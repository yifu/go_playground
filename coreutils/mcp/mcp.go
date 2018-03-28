package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TODO Implement a test suite of some sort.
// TODO implémenter cp -r
// TODO implémenter cp -f
// TODO Implémenter un type spécial: path. Faire remonter ces définition en haut du fichier.

type pathList []string
type environment struct {
	isRec bool
	name  string
}

// TODO Make test with symlink in the middle. Maybe we should use os.Lstat() instead.

func main() {
	env := environment{}

	flag.BoolVar(&env.isRec, "r", false, "Copy directories recursively.")
	flag.Parse()

	setUsage()
	checkArgsCount()
	exitCode := 0

	dst, srcs := parseCmdLine()
	oks := make(pathList, 0, len(srcs))

	var dir string
	if len(srcs) == 1 {
		if info, err := os.Stat(dst); err != nil {
			if os.IsNotExist(err) {
				dir, env.name = filepath.Split(dst)

			} else {
				printErr(err)
				os.Exit(2)
			}
		} else {
			if info.IsDir() {
				dir = dst
			} else {
				dir = filepath.Dir(dst)
			}
		}
	} else if len(srcs) > 1 {
		// When len(srcs) > 1, then dst must be a dir. We check that it is the case:
		if !checkIsDir(dst) {
			printErr(NotADirErr{dst})
			os.Exit(1)
		}
		// TODO When -r is used, we must use filepath.Split() in order to get the dir and the filename. In case of "mcp -r dir1 dir2/" then we must have dir2/dir1 as a result.
		dir = dst
	}

	for _, src := range srcs {
		dst := filepath.Join(dir, filepath.Base(src))
		if oks.contains(dst) {
			printErr(WillNotOverwriteErr{dst, src})
			exitCode = 1
			continue
		}

		var err error
		env, err = cp(dir, src, env)
		if err != nil {
			printErr(err)
			exitCode = 1
			continue
		}

		oks = append(oks, dst)
	}
	os.Exit(exitCode)
}

func parseCmdLine() (string, pathList) {
	args := flag.Args()
	dst := filepath.Clean(args[len(args)-1])
	srcs := args[0 : len(args)-1]
	for i, path := range srcs {
		srcs[i] = filepath.Clean(path)
	}
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
	paramCount := len(flag.Args())
	if paramCount < 2 {
		flag.Usage()
		os.Exit(1)
	}
}

func checkIsDir(dst string) bool {
	dirfi, err := os.Stat(dst)
	return err == nil && dirfi.IsDir()
}

// TODO Remove that function.
func mkDst(dst, src string) string {
	if dstfi, err := os.Stat(dst); err != nil {
		if os.IsNotExist(err) {
			// Nothing to do: dst is a filename,
			// which does not exist yet.
			// The user is asking to copy and rename at the same time.
		} else {
			printErr(err)
			os.Exit(2)
		}
	} else if dstfi.IsDir() {
		return filepath.Join(dst, filepath.Base(src))
	}
	// Nothing to do: dst is a filename, we must copy into it directly.
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

func removePrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

func cp(dir, src string, env environment) (environment, error) {
	srcfi, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return env, NoSuchFileOrDirErr{src}
		}
		printErr(err)
		os.Exit(2)
	}

	if srcfi.IsDir() {
		if env.isRec {

			// TODO Refactor walkFunc into its own func.
			err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					printErr(err)
					os.Exit(2)
				}

				path = filepath.Clean(path)

				fmt.Printf("visited file: %q\n", path)
				if info.IsDir() {
					fmt.Println("path1:", path, ", src:", src)
					// TODO Remove this function to call TrimPrefix() directly.
					path = removePrefix(path, filepath.Dir(src))
					fmt.Println("path2:", path)
					if path == src && env.name != "" {
						path = filepath.Join(filepath.Dir(path), env.name)
						fmt.Println("path3:", path)
					}

					fmt.Println("dir=", dir, ", path=", path)
					// Create a directory with the same name in the destination:
					newdir := filepath.Join(dir, path)
					err := os.Mkdir(newdir, 0700)
					if err != nil {
						if os.IsExist(err) {

						} else {
							return err
						}
					}
				} else {
					fmt.Println("env=", env)
					env, err = cp(filepath.Join(dir, filepath.Dir(path)), path, environment{})
					if err != nil {
						return err
					}
					fmt.Println("env2=", env)
				}

				return nil
			})

			if err != nil {
				printErr(err)
				os.Exit(2)
			}
			// TODO Change back the permission on all the directory files after walking them.
		} else {
			return env, OmittingDirErr{src}
		}
	} else {

		// Open src
		srcf, err := os.OpenFile(src, os.O_RDONLY, 0)
		if err != nil {
			printErr(err)
			os.Exit(2)
		}

		// Stat(dir) first, then only opening it, to avoid creating an empty file.

		// TODO if env.OneSrc && !dstfi.IsDir() { /*Do not touch the destination name*/ } else
		// TODO dst = filepath.Join(dst, filepath.Base(src)) }

		// TODO Rename *fi to *info

		dst := ""
		if env.name != "" {
			dst = filepath.Join(dir, env.name)
		} else {
			dst = filepath.Join(dir, filepath.Base(src))
		}

		if dstfi, err := os.Stat(dst); err != nil {

		} else if os.SameFile(dstfi, srcfi) {
			// Nothing to copy.
			return env, SameFileErr{dst, src}
		}

		dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
		if err != nil {
			printErr(err)
			os.Exit(2)
		}

		if err != nil {
			switch err.(type) {
			case SameFileErr:
				// When src and dst are the same files, there is no copy going on.
				// We just keep going on with the next src.
			default:
				return env, err
			}
		} else {
			if _, err := io.Copy(dstf, srcf); err != nil {
				printErr(err)
				os.Exit(2)
			}
		}
	}
	return env, nil
}

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
	dst, src string
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
