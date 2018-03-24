package main

import (
	"fmt"
)

func producer(i int, c chan int) {
	for {
		c <- i
	}
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)

	go producer(0, c1)
	go producer(1, c2)
	go producer(2, c3)

	for {
		select {
		case i := <-c1:
			fmt.Println(i)
		case i := <-c2:
			fmt.Println(i)
		case i := <-c3:
			fmt.Println(i)
		}
	}
}
