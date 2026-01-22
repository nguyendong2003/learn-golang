package main

import "fmt"

// Định nghĩa Interface Speaker
type Speaker interface {
	Speak() string
}

// Walking là interface rỗng
type Walking interface {
}

// Hearing là interface rỗng
type Hearing any

// Struct Dog
type Dog struct{}

func (d Dog) Speak() string {
	return "Woof!"
}

// Struct Cat
type Cat struct{}

func (c Cat) Speak() string {
	return "Meow!"
}

func makeItSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func makeItWalking(w Walking) {
	fmt.Println(w)
}

func makeItHearing(h Hearing) {
	fmt.Println(h)
}

func main() {
	d := Dog{}
	c := Cat{}

	makeItSpeak(d) // Chạy tốt vì Dog có method Speak()
	makeItSpeak(c) // Chạy tốt vì Cat có method Speak()

	makeItWalking(d) // {}
	makeItWalking(c) // {}

	makeItHearing(d) // {}
	makeItHearing(c) // {}
}
