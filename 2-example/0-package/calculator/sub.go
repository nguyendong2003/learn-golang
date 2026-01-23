package calculator

func Sub(a, b int) int {
	sum := Add(a, -b)
	return sum
}

func SubFloat(a, b float64) float64 {
	return a - b
}
