package main

import "fmt"

func myPointer() {

	var x int32 = 10
	var y bool = true

	pointerX := &x
	pointerY := &y

	var pointerZ *string

	fmt.Printf("%d%d%d",*pointerX,pointerY,pointerZ)

	
}