package main

import (
	"fmt"
	"log"
	"strings"
)

var str = `116
37 23 108 59 86 64 94 14 105 17 111 65 55 31 79 97 78 25 50 22 66 46 104 98 81 90 68 40 103 77 74 18 69 82 41 4 48 83 67 6 2 95 54 100 99 84 34 88 27 72 32 62 9 56 109 115 33 15 91 29 85 114 112 20 26 30 93 96 87 42 38 60 7 73 35 12 10 57 80 13 52 44 16 70 8 39 107 106 63 24 92 45 75 116 5 61 49 101 71 11 53 43 102 110 1 58 36 28 76 47 113 21 89 51 19 3`

func main() {
	r := strings.NewReader(str)
	var w strings.Builder

	n := 0
	if _, err := fmt.Fscanln(r, &n); err != nil {
		log.Panic(err)
	} else {
		fmt.Fprintln(&w, n)
	}

	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(r, &v)
		fmt.Fprintf(&w, "%v", v)
		if i != n-1 {
			fmt.Fprintf(&w, " ")
		}
	}

	fmt.Println("result=", strings.Compare(w.String(), str))
}

// BSTNode is a binary search tree node.
type BSTNode struct {
	value       int
	left, right *BSTNode
}
