package technical

/*
* A Bollinger Band is a techinical analysis volatility indicator developed by John Bollinger in the 1980s.
*
* They provide a relative indication of high and low prices in the market for a given time period.
*
* Bollinger bands consist of a center line (usually a EMA of a stock price) and two price channels (bands)
* above and below the center line scaled by some standard deviation.
 */

// Bound64 represents a pair of lower and upper values and a midpoint
type Bound64 struct {
	Lower    float64 `json:"lower"`
	Midpoint float64 `json:"midpoint"`
	Upper    float64 `json:"upper"`
}

// Bound32 is a 32 bit version of Bound54
type Bound32 struct {
	Lower    float32 `json:"lower"`
	Midpoint float32 `json:"midpoint"`
	Upper    float32 `json:"upper"`
}

// RoundBoundToNearestCent64 ceils a Bollinger Bound upper bound and floors the lower bound to near cent
func RoundBoundToNearestCent64(b *Bound64) {
	// ceil upper
	b.Upper = RoundUp64(b.Upper, 2)

	// floor lower
	b.Lower = RoundDown64(b.Lower, 2)
}

// RoundBoundtoNearestCent32 is 32 bit version of RoundBoundtoNearestCent64
func RoundBoundToNearestCent32(b *Bound32) {
	// ceil upper
	b.Upper = RoundUp32(b.Upper, 2)

	// floor lower
	b.Lower = RoundDown32(b.Lower, 2)
}

// CompareBound64 checks if the first bound is wider than the second bound
func CompareBound64(first *Bound64, second *Bound64) bool {
	firstBoundRange := first.Upper - first.Lower
	secondBoundRange := second.Upper - second.Lower

	if firstBoundRange > secondBoundRange {
		return true
	}

	return false
}

// CompareBound32 is 32 bit version of CompareBound64
func CompareBound32(first *Bound32, second *Bound32) bool {
	firstBoundRange := first.Upper - first.Lower
	secondBoundRange := second.Upper - second.Lower

	if firstBoundRange > secondBoundRange {
		return true
	}

	return false
}

// BollBound64 creates a float64 Bollinger Bound for the given period
// A period is a predetermined segment of a data series
//
// Parameters:
//		period: list of data values
//		k: midpoint of the bound in the given period
//		a (alpha): multiplier on the standard deviation of the period
func BollBound64(period []float64, k float64, a float64) Bound64 {
	var b Bound64

	leg := StdDev64(period) * a

	b.Midpoint = k
	b.Lower = k - leg
	b.Upper = k + leg

	return b
}

// BollBound32 is 32 bit version of BollBound64
func BollBound32(period []float32, k float32, a float32) Bound32 {
	var b Bound32

	leg := StdDev32(period) * a

	b.Midpoint = k
	b.Lower = k - leg
	b.Upper = k + leg

	return b
}

// RollingerBollingerConst64 computes a float64 Bollinger Bound for a given period in a series
// Uses a constant specified midpoint for the Bollinger Bound
// Usage: Call while iterating over a series of values to build a full Bollinger Band.
//
// Parameters:
//		period: list of float values
//		k: static midpoint of the bound for a period
//		a (alpha): multiplier on the standard deviation of the period
func RollingBollingerConst64(period []float64, k float64, a float64) Bound64 {
	return BollBound64(period, k, a)
}

// RollingBollingerConst32 is 32 bit version of RollingBollingerConst64
func RollingBollingerConst32(period []float32, k float32, a float32) Bound32 {
	return BollBound32(period, k, a)
}

// RollingerBollingerSMA64 computes a float64 Bollinger Bound for a given period in a series
// Uses a rolling Simple Moving Average midpoint for the Bollinger Bound
// Usage: Call while iterating over a series of values to build a full Bollinger Band.
//
// Parameters:
//		period: list of float values
//		a (alpha): multiplier on the standard deviation of the period
func RollingBollingerSMA64(period []float64, a float64) Bound64 {
	if period == nil || len(period) == 0 {
		return Bound64{}
	}

	return BollBound64(period, SimpleAvg64(period), a)
}

// RollingBollingerSMA32 is 32 bit version of ROllingBollingerSMA64
func RollingBollingerSMA32(period []float32, a float32) Bound32 {
	if period == nil || len(period) == 0 {
		return Bound32{}
	}

	return BollBound32(period, SimpleAvg32(period), a)
}

// RollingBollingerEMA64 computes a float64 Bollinger Bound for a given period in a series
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
func RollingBollingerEMA64(period []float64, y float64, v float64, last float64, a float64) Bound64 {
	var (
		b Bound64
		k float64
	)

	if period == nil || len(period) == 0 {
		return b
	}

	k = RollingEMA64(v, last, y)

	return BollBound64(period, k, a)
}

// RollingBollingerEMA32 is 32 bit version of RollingBollingerEMA64
func RollingBollingerEMA32(period []float32, y float32, v float32, last float32, a float32) Bound32 {
	var (
		b Bound32
		k float32
	)

	if period == nil || len(period) == 0 {
		return b
	}

	k = RollingEMA32(v, last, y)

	return BollBound32(period, k, a)
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
func StaticBollingerConst64(series []float64, lb int, k float64, a float64) []Bound64 {
	if series == nil || len(series) == 0 {
		return nil
	}

	band := make([]Bound64, len(series))
	for i := range series {
		j := i + 1 // offset by 1 bc of idx
		if j < lb {
			band[i] = Bound64{}
			continue
		}

		band[i] = BollBound64(series[j-lb:j], k, a)
	}

	return band
}

// StaticBollingerConst32 is 32 bit version of StaticBollingerConst64
func StaticBollingerConst32(series []float32, lb int, k float32, a float32) []Bound32 {
	if series == nil || len(series) == 0 {
		return nil
	}

	band := make([]Bound32, len(series))
	for i := range series {
		j := i + 1 // offset by 1 bc of idx
		if j < lb {
			band[i] = Bound32{}
			continue
		}

		band[i] = BollBound32(series[j-lb:j], k, a)
	}

	return band
}

// StaticBollingerSMA64 creates a float64 Bollinger Band using a static standard deviation multiplier and period lookback
// If time series data, assumes ascending order.
// Uses a Simple Moving Average midpoint for the Bollinger Bound
//
// Parameters:
//		series: data series
//		lb: lookback to derive a period
//		a (alpha): multiplier on the standard deviation of the period
func StaticBollingerSMA64(series []float64, lb int, a float64) []Bound64 {
	if series == nil || len(series) == 0 {
		return nil
	}

	band := make([]Bound64, len(series))
	for i := range series {
		j := i + 1 // offset by 1 bc of idx
		if j < lb {
			band[i] = Bound64{}
			continue
		}

		period := series[j-lb : j]

		band[i] = BollBound64(period, SimpleAvg64(period), a)
	}

	return band
}

// StaticBollingerSMA32 is 32 bit version of StaticBollingerSMA64
func StaticBollingerSMA32(series []float32, lb int, a float32) []Bound32 {
	if series == nil || len(series) == 0 {
		return nil
	}

	band := make([]Bound32, len(series))
	for i := range series {
		j := i + 1
		if j < lb {
			band[i] = Bound32{}
			continue
		}

		period := series[j-lb : j]

		band[i] = BollBound32(period, SimpleAvg32(period), a)
	}

	return band
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
func StaticBollingerEMA64(series []float64, lb int, y float64, a float64) []Bound64 {
	if series == nil || len(series) == 0 {
		return nil
	}

	var (
		k    float64
		last float64
	)

	if y == 0.0 { // use default smoothing
		y = 2.0 / float64(lb+1)
	}

	band := make([]Bound64, len(series))
	for i, v := range series {
		j := i + 1
		if j < lb {
			band[i] = Bound64{}
			continue
		}

		period := series[j-lb : j]

		if j == lb { // first bound is simple avg
			k = SimpleAvg64(period)
			last = k
		} else {
			k = RollingEMA64(v, last, y)
			last = k
		}

		band[i] = BollBound64(period, k, a)
	}

	return band
}

// StaticBollingerEMA32 is 32 bit version of StaticBollingerEMA64
func StaticBollingerEMA32(series []float32, lb int, y float32, a float32) []Bound32 {
	if series == nil || len(series) == 0 {
		return nil
	}

	var (
		k    float32
		last float32
	)

	if y == 0.0 { // use default smoothing
		y = 2.0 / float32(lb+1)
	}

	band := make([]Bound32, len(series))
	for i, v := range series {
		j := i + 1 // offset by 1 bc of index
		if j < lb {
			band[i] = Bound32{} // empty bound
			continue
		}

		period := series[j-lb : j] // inclusive of current bound

		if j == lb { // first bound is simple avg
			k = SimpleAvg32(period)
			last = k
		} else {
			k = RollingEMA32(v, last, y)
			last = k
		}

		band[i] = BollBound32(period, k, a)
	}

	return band
}
