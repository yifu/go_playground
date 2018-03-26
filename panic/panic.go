package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello panic.")
	f()
}

func f() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("err=", err)
		}
		panic("NEW ERROR")
	}()
	panic("ERROR")
}