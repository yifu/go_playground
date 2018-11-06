package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// var str = `116
// 37 23 108 59 86 64 94 14 105 17 111 65 55 31 79 97 78 25 50 22 66 46 104 98 81 90 68 40 103 77 74 18 69 82 41 4 48 83 67 6 2 95 54 100 99 84 34 88 27 72 32 62 9 56 109 115 33 15 91 29 85 114 112 20 26 30 93 96 87 42 38 60 7 73 35 12 10 57 80 13 52 44 16 70 8 39 107 106 63 24 92 45 75 116 5 61 49 101 71 11 53 43 102 110 1 58 36 28 76 47 113 21 89 51 19 3`
var str = `6
1 2 5 3 6 4`

func main() {
	r := strings.NewReader(str)
	var w strings.Builder

	n := 0
	if _, err := fmt.Fscanln(r, &n); err != nil {
		log.Panic(err)
	} else {
		fmt.Fprintln(&w, n)
	}

	var root *BSTNode
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(r, &v)
		fmt.Println(v)
		fmt.Fprintf(&w, "%v", v)
		if i != n-1 {
			fmt.Fprintf(&w, " ")
		}
		root = root.insert(v)
	}

	//fmt.Println("cmp =", strings.Compare(w.String(), str))

	// Construct heap
	heap := make([]*BSTNode, 0)
	heap = append(heap, root)
	lastHeapLineWidth := 1
	for !lastLineIsOnlyNil(heap, lastHeapLineWidth) {
		lastLine := heap[len(heap)-lastHeapLineWidth:]
		//fmt.Println("len(heap)=", len(heap))
		//fmt.Println("len(lastLine)=", len(lastLine))
		//fmt.Println("lastHeapLineWidth=", lastHeapLineWidth)
		//fmt.Println("lastHeapLineWidth=", lastHeapLineWidth, "len(heap)-lastHeapLineWidth=", len(heap)-lastNumApp, "len(lastLine)=", len(lastLine))

		for _, p := range lastLine {
			if p == nil {
				heap = append(heap, nil, nil)
			} else {
				heap = append(heap, p.left, p.right)
				fmt.Print("val=", p.value)
			}
		}
		fmt.Println()
		lastHeapLineWidth *= 2
	}
	fmt.Println("lastHeapLineWidth=", lastHeapLineWidth)

	const nodeWidth = 2

	curLvl := 0
	nelt := 0
	nchar := lastHeapLineWidth
	for i, p := range heap {
		lvl := int(math.Log2(float64(i+1))) + 1
		if curLvl != lvl {
			curLvl = lvl
			if nelt == 0 {
				nelt = 1
			} else {
				nelt *= 2
			}
			nchar /= 2
			fmt.Print("|lastHeapLineWidth/nelt=", lastHeapLineWidth/nelt)
			fmt.Println()
		}

		spaces := strings.Repeat(" ", (lastHeapLineWidth/nelt)/2)
		if p == nil {
			fmt.Print(spaces + spaces)
		} else {
			fmt.Print(spaces + strconv.Itoa(p.value) + spaces)
		}
	}

	fmt.Println()
}

func lastLineIsOnlyNil(heap []*BSTNode, lastHeapLineWidth int) bool {
	lastLine := heap[len(heap)-lastHeapLineWidth:]
	for _, p := range lastLine {
		if p != nil {
			return false
		}
	}
	return true
}

// BSTNode is a binary search tree node.
type BSTNode struct {
	value       int
	left, right *BSTNode
}

func (root *BSTNode) insert(v int) *BSTNode {
	if root == nil {
		return &BSTNode{v, nil, nil}
	}
	if v <= root.value {
		root.left = root.left.insert(v)
	} else {
		root.right = root.right.insert(v)
	}
	return root
}
