package main

import (
	"fmt"
	"os"
	"strconv"
)

// Test with: for i in $(seq 1 100); do echo "$i/$(fib $i)";  done

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: fib n")
		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Bad argument: %q\n", os.Args[1])
		os.Exit(1)
	}

	fmt.Printf("%v\n", fib(n))
}

func fib(n int) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return fib(n-1) + fib(n-2)
	}
}
