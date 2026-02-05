package main

// 1. Khai báo Strategy
type PaymentStrategy interface {
	Pay(amount int)
}

// 2. Concrete Strategy (Các chiến lược cụ thể)
type CreditCardPayment struct{}

func (c CreditCardPayment) Pay(amount int) {
	println("Thanh toán", amount, "bằng thẻ tín dụng")
}

type PaypalPayment struct{}

func (p PaypalPayment) Pay(amount int) {
	println("Thanh toán", amount, "bằng PayPal")
}

// 3. Context
type ShoppingCart struct {
	strategy PaymentStrategy
}

func (s *ShoppingCart) SetStrategy(strategy PaymentStrategy) {
	s.strategy = strategy
}

func (s ShoppingCart) Checkout(amount int) {
	if s.strategy == nil {
		println("Chưa chọn phương thức thanh toán")
		return
	}
	s.strategy.Pay(amount)
}

// 4. Sử dụng (Client)
func main() {
	cart := ShoppingCart{}

	cart.SetStrategy(CreditCardPayment{})
	cart.Checkout(100)

	cart.SetStrategy(PaypalPayment{})
	cart.Checkout(200)
}
