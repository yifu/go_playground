package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var inputrdr io.Reader = makeInputReader()
	if _, err := io.Copy(os.Stdout, inputrdr); err != nil {
		fmt.Println("Err:", err)
	}
}

func makeInputReader() io.Reader {
	if len(os.Args) == 1 {
		return os.Stdin
	}

	readers := make([]io.Reader, 0)

	for i, elt := range os.Args {
		if i == 0 {
			continue
		}

		var r io.Reader
		if f, err := os.Open(elt); err != nil {
			r = strings.NewReader(err.Error() + "\n")
		} else {
			r = f
		}

		readers = append(readers, r)
	}

	return io.MultiReader(readers...)
}
