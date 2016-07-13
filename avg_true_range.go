package technical

import (
	"math"
)

/*
Average True Range (ATR) developed by J. Welles Wilder Jr
is a technical analysis indicator that measures the price change volatility.
*/

// TrueRange64 (TR) computes the ATR Wilder True Range for the given period
// True range is defined as the following:
// max(High, last) −  min(Low, last)
// where last is the previous period close, or last value
func TrueRange64(period []float64, last float64) float64 {
	if len(period) == 0 {
		return 0.0
	}

	magic := 4000000000.0 // magic number

	high := 0.0
	low := magic // first val is magic number
	for _, v := range period {
		if v > high {
			high = v
		}

		if v < low && v > 0.0 {
			low = v
		}
	}

	if last == 0.0 {
		v := high - low
		if v < 0.0 {
			return 0.0
		}

		return v
	}

	v := math.Max(high, last) - math.Min(low, last)

	// if period is all 0 or negative return 0
	if v < 0.0 || v == magic {
		return 0.0
	}

	return v
}

func TrueRange32(period []float32, last float32) float32 {
	if len(period) == 0 {
		return 0.0
	}

	magic := float32(4000000000.0) // magic number

	high := float32(0.0)
	low := magic // first val is magic number
	for _, v := range period {
		if v > high {
			high = v
		}

		if v < low && v > 0.0 {
			low = v
		}
	}

	if last == 0.0 {
		v := high - low
		if v < 0.0 {
			return 0.0
		}

		return v
	}

	v := float32(math.Max(float64(high), float64(last)) - math.Min(float64(low), float64(last)))

	// if period is all 0 or negative return 0
	if v < 0.0 || v == magic {
		return 0.0
	}

	return v
}

// RollingATR64 computes the next ATR value based on the prior periods ATR (lastATR),
// the current periods TrueRange64 (curTR), and the number of periods (n)
func RollingATR64(lastATR float64, curTR float64, n int) float64 {
	if n == 0 {
		return 0.0
	}

	return (lastATR*float64(n-1) + curTR) / float64(n)
}

func RollingATR32(lastATR float32, curTR float32, n int) float32 {
	if n == 0 {
		return 0.0
	}

	return (lastATR*float32(n-1) + curTR) / float32(n)
}

// StaticATR64 computes an ATR value from a full length time series based on a
// number of periods (n) and the period size (s)
// It is assumed that the period size (s) is of the same unit of time as
// the indices of the series for the values they represent.
// ex. s = 7 days and a series with 7 elements assumes each index represents
// a value for a daily unit of measurement
// If there are not enough data points in the series to compute a
// complete ATR value the simple average of the partially complete
// true ranges of the derived periods is used as an approximate ATR
// To compute a true range a period must be totally complete
func StaticATR64(series []float64, n int, s int) float64 {
	if n == 0 || s == 0 {
		return 0.0
	}

	trngs := make([]float64, 0)

	// iterate over the series in period size parts and compute true ranges
	for i := 0; i < len(series); i += s {
		if i+s > len(series) {
			break
		}

		period := series[i : i+s]

		// compute and add true range
		last := 0.0
		if i != 0 {
			last = series[i-1]
		}

		trngs = append(trngs, TrueRange64(period, last))

	}

	// if dont have enough period true range values just return avg
	if n >= len(trngs) {
		return SimpleAvg64(trngs)
	}

	// first ATR is simple avg of warm up true ranges
	// i.e. true ranges of the first n complete periods
	// after that, additional ATR values are computed
	// using RollingATR64 providing exponential decay of prior atr
	atr := SimpleAvg64(trngs[0:n])
	for i := n; i < len(trngs); i++ {
		atr = RollingATR64(atr, trngs[i], n)
	}

	return atr
}

func StaticATR32(series []float32, n int, s int) float32 {
	if n == 0 || s == 0 {
		return 0.0
	}

	trngs := make([]float32, 0)

	// iterate over the series in period size parts and compute true ranges
	for i := 0; i < len(series); i += s {
		if i+s > len(series) {
			break
		}

		period := series[i : i+s]

		// compute and add true range
		last := float32(0.0)
		if i != 0 {
			last = series[i-1]
		}

		trngs = append(trngs, TrueRange32(period, last))
	}

	// if dont have enough period true range values just return avg
	if n >= len(trngs) {
		return SimpleAvg32(trngs)
	}

	// first ATR is simple avg of warm up true ranges
	// i.e. true ranges of the first n complete periods
	// after that, additional ATR values are computed
	// using RollingATR32  providing exponential decay of prior atr
	atr := SimpleAvg32(trngs[0:n])
	for i := n; i < len(trngs); i++ {
		atr = RollingATR32(atr, trngs[i], n)
	}

	return atr
}
