package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	fmt.Println("Hello, playground")

	currentWidth := 1
	nextLvlWidth := currentWidth
	curLvl := 1
	for i := 0; i < 32; i++ {
		//fmt.Print("i=", i)
		computedLvl := int(math.Log2(float64(i+1))) + 1
		if curLvl != computedLvl {
			log.Panic(fmt.Sprint("i=", i, ", curLvl=", curLvl, ", computedLvl=", computedLvl))
		}

		nextLvlWidth--
		if nextLvlWidth == 0 {
			currentWidth *= 2
			nextLvlWidth = currentWidth
			curLvl++
			fmt.Println()
		}
	}
}
