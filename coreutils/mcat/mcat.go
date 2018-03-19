package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
	} else {
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

		mr := io.MultiReader(readers...)

		if _, err := io.Copy(os.Stdout, mr); err != nil {
			fmt.Println("Err:", err)
		}
	}
}