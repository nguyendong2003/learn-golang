/*
Ví dụ nhúng Interface vào Struct
*/
package structembedding

import "fmt"

type Brewer interface {
	Brew() string
}

type SmartCoffeeMachine struct {
	Brand  string
	Brewer // Nhúng Interface vào đây
}

// Món 1: Cà phê phin
type VietnameseCoffee struct{}

func (v VietnameseCoffee) Brew() string { return "Cà phê sữa đá đậm đà" }

// Món 2: Trà túi lọc
type TeaBag struct{}

func (t TeaBag) Brew() string { return "Trà nóng thanh tịnh" }

func Main() {
	// Lắp "Cà phê phin" vào máy
	machine := SmartCoffeeMachine{
		Brand:  "GopherBrew",
		Brewer: VietnameseCoffee{},
	}

	// Máy tự động có phương thức Brew() của thứ bên trong
	fmt.Println(machine.Brand, "đang pha:", machine.Brew())

	// Đổi sang "Trà" cực kỳ dễ dàng
	machine.Brewer = TeaBag{}
	fmt.Println(machine.Brand, "đang pha:", machine.Brew())
}
