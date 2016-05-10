package bollinger

import (
	"errors"
	"math"
)

/*
Package for simple calculations of Bollinger Bands.
*/

const ErrEmptyList = "List cannot be empty."
const ErrSmoothingRng = "Custom defined smoothing factor must be less than 1 and greater than 0."

// A Bound represents a pair of lower and upper values and a midpoint.
type Bound64 struct {
	Lower    float64
	MidPoint float64
	Upper    float64
}

type Bound32 struct {
	Lower    float32
	MidPoint float32
	Upper    float32
}

// SimpleAvg computes the simple average of a given list of values
func SimpleAvg64(xs []float64) (float64, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	var sum float64

	for _, v := range xs {
		sum += v
	}

	return sum / float64(len(xs)), nil
}

func SimpleAvg32(xs []float32) (float32, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	var sum float32

	for _, v := range xs {
		sum += v
	}

	return sum / float32(len(xs)), nil
}

// Variance computes the population variance of a given iist of values
func Variance64(xs []float64) (float64, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	var total float64

	avg, err := SimpleAvg64(xs)
	if err != nil {
		return 0.0, err
	}

	for _, v := range xs {
		diff := v - avg
		total += (diff * diff)
	}

	return total / float64(len(xs)), nil
}

func Variance32(xs []float32) (float32, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	var total float32

	avg, err := SimpleAvg32(xs)
	if err != nil {
		return 0.0, err
	}

	for _, v := range xs {
		diff := v - avg
		total += (diff * diff)
	}

	return total / float32(len(xs)), nil
}

// StdDev computes the standard deviation of a given list of values
func StdDev64(xs []float64) (float64, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	v, err := Variance64(xs)
	if err != nil {
		return 0.0, err
	}

	res := math.Sqrt(v)
	if math.IsNaN(res) {
		return 0.0, nil
	} else {
		return res, nil
	}
}

func StdDev32(xs []float32) (float32, error) {
	if xs == nil || len(xs) == 0 {
		return 0.0, errors.New(ErrEmptyList)
	}

	v, err := Variance32(xs)
	if err != nil {
		return 0.0, err
	}

	res := math.Sqrt(float64(v))
	if math.IsNaN(res) {
		return float32(0.0), nil
	} else {
		return float32(res), nil
	}
}

// CompareBound32 checks if the first bound is wider than the second bound
func CompareBound32(first *Bound32, second *Bound32) bool {
	firstBoundRange := first.Upper - first.Lower
	secondBoundRange := second.Upper - second.Lower

	if firstBoundRange > secondBoundRange {
		return true
	} else {
		return false
	}
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
func EwmaSeries64(series []float64, y float64, lb int) ([]float64, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	size := len(series)
	if lb >= size {
		return nil, errors.New("Lookback must be less than the length of the series.")
	}

	if y == 0.0 { // use default smoothing
		y = 2.0 / float64(lb+1)
	}

	var lastEma float64
	ewmas := make([]float64, size)
	for i, v := range series {
		if i < lb {
			ewmas[i] = 0.0
			continue
		} else if i == lb { // first is simple average
			savg, err := SimpleAvg64(series[i-lb : i])
			if err != nil {
				return nil, err
			}

			ewmas[i] = savg
			lastEma = savg
			continue
		} else { // compute ewma
			curEma := RollingEMA64(v, lastEma, y)
			ewmas[i] = curEma
			lastEma = curEma
		}
	}

	return ewmas, nil
}

func EwmaSeries32(series []float32, y float32, lb int) ([]float32, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	size := len(series)
	if lb >= size {
		return nil, errors.New("Lookback must be less than the length of the series.")
	}

	if y == float32(0.0) { // use default smoothing
		y = 2.0 / float32(lb+1)
	}

	var lastEma float32
	ewmas := make([]float32, size)
	for i, v := range series {
		if i < lb {
			ewmas[i] = float32(0.0)
			continue
		} else if i == lb { // first is simple average
			savg, err := SimpleAvg32(series[i-lb : i])
			if err != nil {
				return nil, err
			}

			ewmas[i] = savg
			lastEma = savg
			continue
		} else { // compute ewma
			curEma := RollingEMA32(v, lastEma, y)
			ewmas[i] = curEma
			lastEma = curEma
		}
	}

	return ewmas, nil
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

// BollBound64 creates a float64 Bollinger Bound for the given period
// A period is a predetermined segment of a data series
//
// Parameters:
//		period: list of data values
//		k: midpoint of the bound in the given period
//		a (alpha): multiplier on the standard deviation of the period
func BollBound64(period []float64, k float64, a float64) (Bound64, error) {
	if period == nil || len(period) == 0 {
		return Bound64{}, errors.New(ErrEmptyList)
	}

	var b Bound64

	stdv, err := StdDev64(period)
	if err != nil {
		return Bound64{}, err
	}

	leg := stdv * a

	b.MidPoint = k
	b.Lower = k - leg
	b.Upper = k + leg

	return b, nil
}

func BollBound32(period []float32, k float32, a float32) (Bound32, error) {
	if period == nil || len(period) == 0 {
		return Bound32{}, errors.New(ErrEmptyList)
	}

	var b Bound32

	stdv, err := StdDev32(period)
	if err != nil {
		return Bound32{}, err
	}

	leg := stdv * a

	b.MidPoint = k
	b.Lower = k - leg
	b.Upper = k + leg

	return b, nil
}

// RollingerBollingerConst64 computes a float64 Bollinger Bound for a given period in a series
// Uses a constant specified midpoint for the Bollinger Bound
// Usage: Call while iterating over a series of values to build a full Bollinger Band.
//
// Parameters:
//		period: list of float values
//		k: static midpoint of the bound for a period
//		a (alpha): multiplier on the standard deviation of the period
func RollingBollingerConst64(period []float64, k float64, a float64) (Bound64, error) {
	if period == nil || len(period) == 0 {
		return Bound64{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound64
	)

	b, err = BollBound64(period, k, a)
	if err != nil {
		return Bound64{}, err
	} else {
		return b, nil
	}
}

func RollingBollingerConst32(period []float32, k float32, a float32) (Bound32, error) {
	if period == nil || len(period) == 0 {
		return Bound32{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound32
	)

	b, err = BollBound32(period, k, a)
	if err != nil {
		return Bound32{}, err
	} else {
		return b, nil
	}
}

// RollingerBollingerSMA64 computes a float64 Bollinger Bound for a given period in a series
// Uses a rolling Simple Moving Average midpoint for the Bollinger Bound
// Usage: Call while iterating over a series of values to build a full Bollinger Band.
//
// Parameters:
//		period: list of float values
//		a (alpha): multiplier on the standard deviation of the period
func RollingBollingerSMA64(period []float64, a float64) (Bound64, error) {
	if period == nil || len(period) == 0 {
		return Bound64{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		k   float64
		b   Bound64
	)

	k, err = SimpleAvg64(period)
	if err != nil {
		return Bound64{}, err
	}

	b, err = BollBound64(period, k, a)
	if err != nil {
		return Bound64{}, err
	} else {
		return b, nil
	}
}

func RollingBollingerSMA32(period []float32, a float32) (Bound32, error) {
	if period == nil || len(period) == 0 {
		return Bound32{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		k   float32
		b   Bound32
	)

	k, err = SimpleAvg32(period)
	if err != nil {
		return Bound32{}, err
	}

	b, err = BollBound32(period, k, a)
	if err != nil {
		return Bound32{}, err
	} else {
		return b, nil
	}
}

// RollingerBollingerEMA64 computes a float64 Bollinger Bound for a given period in a series
// Uses a rolling Exponentially Weighted Moving Average midpoint for the Bollinger Bound
// Usage: Call while iterating over a series of values to build a full Bollinger Band.
//
// Parameters:
//		period: list of float values
//		y (lambda): smoothing factor for EWMA
//		v: current underlying value in the series
//		last: last EMA midpoint of the prior bound. The first bound midpoint should be computed using a Simple Average
//		a (alpha): multiplier on the standard deviation of the period
// CONSTAINT: 0.0 < y < 1.0
func RollingBollingerEMA64(period []float64, y float64, v float64, last float64, a float64) (Bound64, error) {
	if period == nil || len(period) == 0 {
		return Bound64{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound64
		k   float64
	)

	k = RollingEMA64(v, last, y)

	b, err = BollBound64(period, k, a)
	if err != nil {
		return Bound64{}, err
	} else {
		return b, nil
	}
}

func RollingBollingerEMA32(period []float32, y float32, v float32, last float32, a float32) (Bound32, error) {
	if period == nil || len(period) == 0 {
		return Bound32{}, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound32
		k   float32
	)

	k = RollingEMA32(v, last, y)

	b, err = BollBound32(period, k, a)
	if err != nil {
		return Bound32{}, err
	} else {
		return b, nil
	}
}

// StaticBollingerConst64 creates a float64 Bollinger Band using a static standard deviation multiplier and period lookback
// If time series data, assumes ascending order.
// Uses a constant specified midpoint for the Bollinger Bound
//
// Parameters:
//		series: data series
//		lb: lookback to derive a period
//		k: static midpoint of the bound for a period
//		a (alpha): multiplier on the standard deviation of the period
func StaticBollingerConst64(series []float64, lb int, k float64, a float64) ([]Bound64, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound64
	)

	band := make([]Bound64, len(series))
	for i, _ := range series {
		if i < lb {
			band[i] = Bound64{}
			continue
		}

		b, err = BollBound64(series[i-lb:i], k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}

func StaticBollingerConst32(series []float32, lb int, k float32, a float32) ([]Bound32, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound32
	)

	band := make([]Bound32, len(series))
	for i, _ := range series {
		if i < lb {
			band[i] = Bound32{}
			continue
		}

		b, err = BollBound32(series[i-lb:i], k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}

// StaticBollingerSMA64 creates a float64 Bollinger Band using a static standard deviation multiplier and period lookback
// If time series data, assumes ascending order.
// Uses a Simple Moving Average midpoint for the Bollinger Bound
//
// Parameters:
//		series: data series
//		lb: lookback to derive a period
//		a (alpha): multiplier on the standard deviation of the period
func StaticBollingerSMA64(series []float64, lb int, a float64) ([]Bound64, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound64
		k   float64
	)

	band := make([]Bound64, len(series))
	for i, _ := range series {
		if i < lb {
			band[i] = Bound64{}
			continue
		}

		period := series[i-lb : i]

		k, err = SimpleAvg64(period)
		if err != nil {
			return nil, err
		}

		b, err = BollBound64(period, k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}

func StaticBollingerSMA32(series []float32, lb int, a float32) ([]Bound32, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err error
		b   Bound32
		k   float32
	)

	band := make([]Bound32, len(series))
	for i, _ := range series {
		if i < lb {
			band[i] = Bound32{}
			continue
		}

		period := series[i-lb : i]

		k, err = SimpleAvg32(period)
		if err != nil {
			return nil, err
		}

		b, err = BollBound32(period, k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}

// StaticBollingerEMA64 creates a float64 Bollinger Band using a static standard deviation multiplier and period lookback
// If time series data, assumes ascending order.
// Uses an Exponentially Weighted Moving Average midpoint for the Bollinger Bound
//
// Parameters:
//		series: data series
//		lb: lookback to derive a period
//		y (lambda): smoothing factor for rolling EWMA. If 0.0 use default formulaic calculation.
//		a (alpha): multiplier on the standard deviation of the period
func StaticBollingerEMA64(series []float64, lb int, y float64, a float64) ([]Bound64, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err  error
		b    Bound64
		k    float64
		last float64
	)

	if y == 0.0 { // use default smoothing
		y = 2.0 / float64(lb+1)
	}

	band := make([]Bound64, len(series))
	for i, v := range series {
		if i < lb {
			band[i] = Bound64{}
			continue
		}

		period := series[i-lb : i]

		if i == lb { // first bound is simple avg
			k, err = SimpleAvg64(period)
			if err != nil {
				return nil, err
			}
			last = k
		} else {
			k = RollingEMA64(v, last, y)
			last = k
		}

		b, err = BollBound64(period, k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}

func StaticBollingerEMA32(series []float32, lb int, y float32, a float32) ([]Bound32, error) {
	if series == nil || len(series) == 0 {
		return nil, errors.New(ErrEmptyList)
	}

	var (
		err  error
		b    Bound32
		k    float32
		last float32
	)

	if y == 0.0 { // use default smoothing
		y = 2.0 / float32(lb+1)
	}

	band := make([]Bound32, len(series))
	for i, v := range series {
		if i < lb {
			band[i] = Bound32{} // empty bound
			continue
		}

		period := series[i-lb : i]

		if i == lb { // first bound is simple avg
			k, err = SimpleAvg32(period)
			if err != nil {
				return nil, err
			}
			last = k
		} else {
			k = RollingEMA32(v, last, y)
			last = k
		}

		b, err = BollBound32(period, k, a)
		if err != nil {
			return nil, err
		}

		band[i] = b
	}

	return band, nil
}
