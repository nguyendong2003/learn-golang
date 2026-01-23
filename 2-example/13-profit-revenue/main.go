package profitrevenue

import (
	"fmt"
)

func Main() {
	var revenue int
	var expenses int
	var taxRate float64

	fmt.Println("Enter total revenue:")
	fmt.Scanln(&revenue)

	fmt.Println("Enter total expenses:")
	fmt.Scanln(&expenses)

	fmt.Println("Enter tax rate (as a decimal):")
	fmt.Scanln(&taxRate)

	ebt := revenue - expenses
	profit := float64(ebt) * (1 - taxRate/100)
	ratio := float64(ebt) / profit

	fmt.Printf("Profit after tax: %.2f\n", profit)
	fmt.Printf("EBT to Profit Ratio: %.2f\n", ratio)
}
