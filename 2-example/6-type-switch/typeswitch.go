package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}
func check(i interface{}) {
	switch i.(type) {
	case nil:
		fmt.Println("i is nil")
	case int, int32, int64:
		fmt.Println("i is an integer type")
	case string:
		fmt.Println("i is a string")
	default:
		fmt.Println("i is of a different type")
	}
}

func main() {
	do(21)
	do("hello")
	do(true)

	check(nil)
	check(10)
	check("test")
	check(3.14)
}
