package main

import (
	"fmt"
)

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Println("arr len=", len(arr), ", cap=", cap(arr))
	fmt.Printf("arr=||%p||\n", &arr)
	for _, elt := range arr {
		fmt.Println(elt)
	}

	slice := arr[0:len(arr)]
	slice[0] = 100

	fmt.Println("slice, len=", len(slice), ", cap=", cap(slice))
	fmt.Printf("slice=%p\n", slice)
	for _, elt := range slice {
		fmt.Println(elt)
	}

	slice = append(slice, 6, 7, 8)
	slice[0] = 200

	fmt.Println("slice2 len=", len(slice), ", cap=", cap(slice))
	fmt.Printf("slice=%p\n", slice)
	for _, elt := range slice {
		fmt.Println(elt)
	}

	fmt.Println("arr2 len=", len(arr), ", cap=", cap(arr))
	fmt.Printf("arr=||%p||\n", &arr)
	for _, elt := range arr {
		fmt.Println(elt)
	}

}
