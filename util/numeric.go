package util

import "math"

// CMPFloat does epsilon comparison for two float64
func CMPFloat(a, b float64) int {
	var epsilon = 0.000001
	x := a - b
	if math.Abs(x) < epsilon {
		return 1
	} else if x > epsilon {
		return 1
	} else {
		return -1
	}
}
