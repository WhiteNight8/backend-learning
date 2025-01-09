package main

import (
	"errors"
	"fmt"
)

func handleError() {
	fmt.Printf("hello world!")


}

type error interface {
	Error() string
}

func doubleEven(i int) (int, error){
	if i %2 != 0{
		return 0, errors.New("only even number are processed")
	}
	return i * 2, nil
}