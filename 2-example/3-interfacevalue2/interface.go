package main

import "fmt"

type I interface {
	M()
}

type T struct {
	S string
}

func (t *T) M() {
	// Nếu không có kiểm tra nil này, chương trình sẽ panic dòng fmt.Println(t.S) khi t là nil
	if t == nil {
		fmt.Println("<<nil>>")
		return
	}
	fmt.Println(t.S)
}

func main() {
	var i I
	describe(i) // (<nil>, <nil>)

	var t *T
	describe(t) // (*main.T, <nil>)
	i = t
	describe(i) // (*main.T, <nil>)
	i.M()       // <<nil>>

	i = &T{"hello"}
	describe(i) // (*main.T, &{hello})
	i.M()       // hello
}

func describe(i I) {
	fmt.Printf("(%T, %v)\n", i, i)
}
