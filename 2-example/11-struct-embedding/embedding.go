/*
Ví dụ nhúng Struct vào Struct
*/
package structembedding

import "fmt"

type Person struct {
	Name    string
	Age     int
	Address string
}

func (p Person) Greet() {
	fmt.Printf("Hi, I'm %s\n", p.Name)
}

func (p Person) Walk() {
	fmt.Printf("I'm walking. I'm %d \n", p.Age)
}

func (p Person) GetAddress() string {
	return p.Address
}

type Employee struct {
	Person  // Đây là Struct Embedding
	ID      int
	Address string
}

func (e Employee) GetAddress() string {
	return e.Address
}

func ExampleEmbedding() {
	e := Employee{
		Person:  Person{Name: "An", Age: 30, Address: "123 Main St"},
		ID:      101,
		Address: "456 Corporate Blvd",
	}

	fmt.Println(e.Name) // Truy cập trực tiếp thay vì e.Person.Name
	fmt.Println(e.Age)  // Truy cập trực tiếp thay vì e.Person.Age
	fmt.Println(e.ID)
	fmt.Println(e.Address)        // Lưu ý: e.Address sẽ truy cập trường Address của Employee, không phải của Person
	fmt.Println(e.Person.Address) // Truy cập trường Address của Person

	e.Greet()                          // Gọi trực tiếp phương thức của Person
	e.Walk()                           // Gọi trực tiếp phương thức của Person
	fmt.Println(e.GetAddress())        // Gọi phương thức GetAddress của Employee
	fmt.Println(e.Person.GetAddress()) // Gọi phương thức GetAddress của Person
}
