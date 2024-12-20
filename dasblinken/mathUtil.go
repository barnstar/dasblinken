package dasblinken

// Subs two floats and clamps the result at 0.0
func Fsub_norm(a, b float64) float64 {
	diff := a - b
	diff = max(diff, 0.0)
	return diff
}

// Adds two floats and clamps the result at 1.0
func Fadd_norm(a, b float64) float64 {
	sum := a + b
	sum = min(sum, 1.0)
	return sum
}
