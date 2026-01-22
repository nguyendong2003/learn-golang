package main

import "fmt"

func main() {
	var i interface{}
	describe(i) // (<nil>, <nil>)

	i = 42
	describe(i) // (int, 42)

	i = "hello"
	describe(i) // (string, hello)
}

func describe(i interface{}) {
	fmt.Printf("(%T, %v)\n", i, i)
}
