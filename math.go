package assocentity

import "math"

// Average from float slice
func avg(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}

	return total / float64(len(xs))
}

// Round two decimal places
func round(x float64) float64 {
	return math.Round(x*100) / 100
}
