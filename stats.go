package technical

import (
	"math"
)

// RoundUp64 rounds the given number up to the nearest nth decimal place
func RoundUp64(x float64, n int) float64 {
	roundFactor := math.Pow10(n)

	return math.Ceil((x * roundFactor)) / roundFactor
}

func RoundUp32(x float32, n int) float32 {
	roundFactor := math.Pow10(n)

	return float32(math.Ceil((float64(x) * roundFactor)) / roundFactor)
}

// RoundDown64 rounds the given number down to the nearest nth decimal place
func RoundDown64(x float64, n int) float64 {
	roundFactor := math.Pow10(n)

	return math.Floor((x * roundFactor)) / roundFactor
}

func RoundDown32(x float32, n int) float32 {
	roundFactor := math.Pow10(n)

	return float32(math.Floor((float64(x) * roundFactor)) / roundFactor)
}

// SimpleAvg64 computes the simple average of a given list of values
func SimpleAvg64(xs []float64) float64 {
	if len(xs) == 0 {
		return 0.0
	}

	var sum float64

	for _, v := range xs {
		sum += v
	}

	return sum / float64(len(xs))
}

func SimpleAvg32(xs []float32) float32 {
	if len(xs) == 0 {
		return 0.0
	}

	var sum float32

	for _, v := range xs {
		sum += v
	}

	return sum / float32(len(xs))
}

// Variance64 computes the population variance of a given iist of values
func Variance64(xs []float64) float64 {
	if len(xs) == 0 {
		return 0.0
	}

	var total float64

	avg := SimpleAvg64(xs)

	for _, v := range xs {
		diff := v - avg
		total += (diff * diff)
	}

	return total / float64(len(xs))
}

func Variance32(xs []float32) float32 {
	if len(xs) == 0 {
		return 0.0
	}

	var total float32

	avg := SimpleAvg32(xs)

	for _, v := range xs {
		diff := v - avg
		total += (diff * diff)
	}

	return total / float32(len(xs))
}

// StdDev64 computes the standard deviation of a given list of values
func StdDev64(xs []float64) float64 {
	if len(xs) == 0 {
		return 0.0
	}

	res := math.Sqrt(Variance64(xs))

	if math.IsNaN(res) {
		return 0.0
	}

	return res
}

func StdDev32(xs []float32) float32 {
	if xs == nil || len(xs) == 0 {
		return 0.0
	}

	res := math.Sqrt(float64(Variance32(xs)))

	if math.IsNaN(res) {
		return float32(0.0)
	}

	return float32(res)
}

// EwmaSeries computes a list of Exponentially Weighted Moving Averages for a given list of values
// If xs is time series data assumes ascending time order
//
// Parameters:
//		series: the data series
// 		smoothingMethod: 0 if use default formula. 1 if custom
//		y (lambda): decay smoothing factor on the weight of each element
//		lb (lookback): size of the period to compute avg. Must be < len of data series
//
// Constraint: 0 < y < 1
//
func EwmaSeries64(series []float64, y float64, lb int) []float64 {
	if len(series) == 0 {
		return nil
	}

	size := len(series)
	if lb > size { // use full series
		lb = size
	}

	if y == 0.0 { // use default smoothing
		y = 2.0 / float64(lb+1)
	}

	var lastEma float64
	ewmas := make([]float64, size)
	for i, v := range series {
		j := i + 1 // offset 1 bc of idx
		switch {
		case j < lb:
			ewmas[i] = 0.0
		case j == lb: // first is simple average
			savg := SimpleAvg64(series[j-lb : j])
			ewmas[i] = savg
			lastEma = savg
		default: // compute ewma
			curEma := RollingEMA64(v, lastEma, y)
			ewmas[i] = curEma
			lastEma = curEma
		}
	}

	return ewmas
}

func EwmaSeries32(series []float32, y float32, lb int) []float32 {
	if series == nil || len(series) == 0 {
		return nil
	}

	size := len(series)
	if lb > size { // use full series
		lb = size
	}

	if y == float32(0.0) { // use default smoothing
		y = 2.0 / float32(lb+1)
	}

	var lastEma float32
	ewmas := make([]float32, size)
	for i, v := range series {
		j := i + 1 // offset by 1 bc of index
		switch {
		case j < lb:
			ewmas[i] = float32(0.0)
		case j == lb: // first is simple average
			savg := SimpleAvg32(series[j-lb : j])
			ewmas[i] = savg
			lastEma = savg
		default: // compute ewma
			curEma := RollingEMA32(v, lastEma, y)
			ewmas[i] = curEma
			lastEma = curEma
		}
	}

	return ewmas
}

// RollingEwma computes the next EWMA value in a series
// Assumes last EWMA value given is correct
//
// Parameters:
//		v: the current value in the series
//		last: the last EWMA value of the series
//		y (lambda): smoothing factor
// Constraint: 0 < y < 1
func RollingEMA64(v float64, last float64, y float64) float64 {
	return v*y + (1-y)*last
}

func RollingEMA32(v float32, last float32, y float32) float32 {
	return v*y + (1-y)*last
}
