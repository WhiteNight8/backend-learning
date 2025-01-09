package main

func example() {
	defer func() int {
		return 2
	}()
}
