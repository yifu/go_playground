package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if _, err := io.Copy(os.Stdout, makeInputReader()); err != nil {
		fmt.Println("Err:", err)
	}
}

func makeInputReader() io.Reader {
	if len(os.Args) == 1 {
		return os.Stdin
	}

	files := make([]io.Reader, 0)

	for i, filename := range os.Args {
		if i == 0 {
			continue
		}
		files = append(files, openFile(filename))
	}

	return io.MultiReader(files...)
}

func openFile(filename string) io.Reader {
	if f, err := os.Open(filename); err != nil {
		return strings.NewReader(err.Error() + "\n")
	} else {
		return f
	}
}
