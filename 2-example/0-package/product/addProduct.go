package product

import (
	"fmt"
	"learngo/0-package/calculator"
)

func AddProduct(a, b int) int {
	return calculator.Add(a, b)
}

func SubProduct(a, b int) int {
	return calculator.Sub(a, b)
}

func RunProduct() {
	fmt.Println("Sum of products:", AddProduct(10, 5))
	fmt.Println("Difference of products:", SubProduct(10, 4))
}
