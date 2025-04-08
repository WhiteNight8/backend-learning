package main

import "fmt"

var x [3]int

var cx = [12]int{1, 5: 4, 6, 10: 100, 15}

var cxx = []int{1, 2, 4}

func test() {
	var x []int
	x = append(x, 10)

	var s string = "hello"
	fmt.Println(len(s))

	var nilMap map[string]int
	fmt.Println(nilMap)

	ages := make(map[int]string, 10)
	fmt.Println(ages)

	type person struct {
		name string
		age  int
		pet  string
	}

	type people struct {
	}

}
