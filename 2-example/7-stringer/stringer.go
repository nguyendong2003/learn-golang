package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

type PersonStringer struct {
	Name    string
	Address string
	Age     int
}

func (ps PersonStringer) String() string {
	return fmt.Sprintf("%s sống tại %s, %d tuổi", ps.Name, ps.Address, ps.Age)
}

func main() {
	p := Person{"Hoàng", 25}
	fmt.Println(p) // Kết quả: {Hoàng 25}

	ps := PersonStringer{"Hoàng", "Hà Nội", 25}
	fmt.Println(ps) // Kết quả: Hoàng sống tại Hà Nội, 25 tuổi
}
