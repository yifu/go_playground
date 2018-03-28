package main

import (
	"fmt"
	"os"
)

func main() {
	if _, err := os.OpenFile("./ldksjfklsjdfkljsdf.txt", os.O_RDONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
	} else {
		fmt.Println("ok")
	}

	fmt.Println("====")

	if _, err := os.OpenFile("./file.go/toto.txt", os.O_RDONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
	} else {
		fmt.Println("ok")
	}

	fmt.Println("====")

	if f, err := os.OpenFile("./file.go", os.O_RDONLY, 0); err != nil {
		fmt.Println("file.go no ok", err)
		fmt.Printf("type = '%T'\n", err)
	} else {
		fmt.Println("file.go ok.")
		if fstat, err := f.Stat(); err != nil {
			
		} else {
			fmt.Printf("name = %q\n", fstat.Name())	
		}
	}

	fmt.Println("====")

	if dir, err := os.OpenFile("./tmpdir", os.O_RDONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
	} else {
		fmt.Println("ok, ")
		defer dir.Close()

		if dirstat, err := dir.Stat(); err != nil {
			fmt.Println("dir stat, err=", err)
		} else {
			fmt.Println("is dir = ", dirstat.IsDir())
		}
		
		if f, err := os.OpenFile("./tmpdir/toto.txt", os.O_WRONLY|os.O_CREATE, 0); err != nil {
			fmt.Println(err)
		} else {
			defer f.Close()
			f.WriteString("hello world\n")
		}
	}


	fmt.Println("====")

	if dir, err := os.OpenFile("./tmpdir", os.O_WRONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
		if v, ok := err.(*os.PathError); ok {
			fmt.Printf("Path erro = %#v\n", v)
		}
	} else {
		fmt.Println("ok, ")
		defer dir.Close()
	}


	fmt.Println("====")

	if dir, err := os.OpenFile("./tmpdir", os.O_RDWR, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
		if v, ok := err.(*os.PathError); ok {
			fmt.Printf("Path erro = %#v\n", v)
		}
	} else {
		fmt.Println("ok, ")
		defer dir.Close()
	}



	fmt.Println("====")

	if dir, err := os.OpenFile("./tmpdir", os.O_WRONLY, 0); err != nil {
		fmt.Println("no ok", err)
		fmt.Printf("type = '%T'\n", err)
		if v, ok := err.(*os.PathError); ok {
			fmt.Printf("Path erro = %#v\n", v)
		}
	} else {
		fmt.Println("ok, ")
		defer dir.Close()
	}
}