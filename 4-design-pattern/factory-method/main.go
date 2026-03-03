package main

// 1. Product interface
type Transport interface {
	Deliver()
}

// 2. Concrete Products
type Truck struct{}

func (t Truck) Deliver() {
	println("Giao hàng bằng xe tải")
}

type Ship struct{}

func (s Ship) Deliver() {
	println("Giao hàng bằng tàu thủy")
}

// 3. Creator interface
type Logistics interface {
	CreateTransport() Transport // Factory Method
	PlanDelivery()
}

// 4. Concrete Creators
type RoadLogistics struct{}

func (r RoadLogistics) CreateTransport() Transport {
	return Truck{}
}

func (r RoadLogistics) PlanDelivery() {
	transport := r.CreateTransport()
	transport.Deliver()
}

type SeaLogistics struct{}

func (s SeaLogistics) CreateTransport() Transport {
	return Ship{}
}

func (s SeaLogistics) PlanDelivery() {
	transport := s.CreateTransport()
	transport.Deliver()
}

// 5. Client dùng (👉 Client không cần biết Truck hay Ship được tạo ra như thế nào.)
func main() {
	var logistics Logistics

	logistics = RoadLogistics{}
	logistics.PlanDelivery()

	logistics = SeaLogistics{}
	logistics.PlanDelivery()
}
