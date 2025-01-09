package main

import "fmt"

func main1() {
	x := 10

	if x > 5 {
		fmt.Println(x)
		x := 5
		fmt.Println(x)
	}
	fmt.Println(x)
}
