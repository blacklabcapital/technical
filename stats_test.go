package technical

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat"
)

func TestRoundUp64(t *testing.T) {
	assert.Equal(t, 10.1235, RoundUp64(10.12347, 4))
	assert.Equal(t, 10.1234, RoundUp64(10.1234, 4))
	assert.Equal(t, 10.124, RoundUp64(10.1234, 3))
	assert.Equal(t, 10.1, RoundUp64(10.1, 3))
}

func TestRoundUp32(t *testing.T) {
	assert.Equal(t, float32(10.1235), RoundUp32(float32(10.12347), 4))
	assert.Equal(t, float32(10.1234), RoundUp32(float32(10.1234), 4))
	assert.Equal(t, float32(10.124), RoundUp32(float32(10.1234), 3))
	assert.Equal(t, float32(10.101), RoundUp32(float32(10.1), 3))
}

func TestRoundDown64(t *testing.T) {
	assert.Equal(t, 10.1234, RoundDown64(10.12347, 4))
	assert.Equal(t, 10.1234, RoundDown64(10.1234, 4))
	assert.Equal(t, 10.123, RoundDown64(10.1234, 3))
	assert.Equal(t, 10.1, RoundDown64(10.1, 3))
}

func TestRoundDown32(t *testing.T) {
	assert.Equal(t, float32(10.1234), RoundDown32(float32(10.12347), 4))
	assert.Equal(t, float32(10.1233), RoundDown32(float32(10.1234), 4))
	assert.Equal(t, float32(10.123), RoundDown32(float32(10.1234), 3))
	assert.Equal(t, float32(10.1), RoundDown32(float32(10.1), 3))
}

func TestSimpleAvg64(t *testing.T) {
	var v float64

	// nil list
	var l []float64
	v = SimpleAvg64(l)
	assert.Equal(t, 0.0, v)

	// empty list
	l = make([]float64, 0)
	v = SimpleAvg64(l)
	assert.Equal(t, 0.0, v)

	l = append(l, 2.0)
	l = append(l, 4.0)

	// good
	v = SimpleAvg64(l)
	assert.Equal(t, 3.0, v)
}

func TestSimpleAvg32(t *testing.T) {
	var v float32

	// nil list
	var l []float32
	v = SimpleAvg32(l)
	assert.Equal(t, float32(0.0), v)

	// empty list
	l = make([]float32, 0)
	v = SimpleAvg32(l)
	assert.Equal(t, float32(0.0), v)

	l = append(l, 2.0)
	l = append(l, 4.0)

	// good
	v = SimpleAvg32(l)
	assert.Equal(t, float32(3.0), v)
}

func TestVariance64(t *testing.T) {
	var v float64

	// nil list
	var l []float64
	v = Variance64(l)
	assert.Equal(t, 0.0, v)

	// empty list
	l = make([]float64, 0)
	v = Variance64(l)
	assert.Equal(t, 0.0, v)

	// good
	l = []float64{
		1,
		3,
		5,
		7,
	}

	v = Variance64(l)
	assert.Equal(t, 5.0, v)
}

func TestVariance32(t *testing.T) {
	var v float32

	// nil list
	var l []float32
	v = Variance32(l)
	assert.Equal(t, float32(0.0), v)

	// empty list
	l = make([]float32, 0)
	v = Variance32(l)
	assert.Equal(t, float32(0.0), v)

	// good
	l = []float32{
		1,
		3,
		5,
		7,
	}

	v = Variance32(l)
	assert.Equal(t, float32(5.0), v)
}

func TestStdDev64(t *testing.T) {
	var v float64

	// nil list
	var l []float64
	v = StdDev64(l)
	assert.Equal(t, 0.0, v)

	// empty list
	l = make([]float64, 0)
	v = StdDev64(l)
	assert.Equal(t, 0.0, v)

	// good
	l = []float64{
		10,
		100,
	}

	v = StdDev64(l)
	assert.Equal(t, 45.0, v)
}

func TestStdDev32(t *testing.T) {
	var v float32

	// nil list
	var l []float32
	v = StdDev32(l)
	assert.Equal(t, float32(0.0), v)

	// empty list
	l = make([]float32, 0)
	v = StdDev32(l)
	assert.Equal(t, float32(0.0), v)

	// good
	l = []float32{
		10,
		100,
	}

	v = StdDev32(l)
	assert.Equal(t, float32(45), v)
}

func TestEwmaSeries64(t *testing.T) {
	var ewmas []float64

	// nil list
	var l []float64
	ewmas = EwmaSeries64(l, 0.0, 10)
	assert.Nil(t, ewmas)

	// empty list
	l = make([]float64, 0)
	ewmas = EwmaSeries64(l, 0.0, 10)
	assert.Nil(t, ewmas)

	// list < lookback
	l = []float64{
		1.0,
		9.0,
	}

	// should be = sma
	ewmas = EwmaSeries64(l, 0.0, 10)
	assert.Equal(t, 5.0, ewmas[len(ewmas)-1])

	// full data series on lb = 600, decay = default
	// read in test series file
	testseries := make([]float64, 0)
	f, err := os.Open("./mock/test_series.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		num, _ := strconv.ParseFloat(scanner.Text(), 64)
		testseries = append(testseries, num)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	ewmas = EwmaSeries64(testseries, 0.0, 600)
	assert.Equal(t, 39.855408877742725, ewmas[len(ewmas)-1])

	// full data series on lb = 1200, decay = default
	ewmas = EwmaSeries64(testseries, 0.0, 1200)
	assert.Equal(t, 39.9488124468039, ewmas[len(ewmas)-1])

}

func TestEwmaSeries32(t *testing.T) {
	var ewmas []float32

	// nil list
	var l []float32
	ewmas = EwmaSeries32(l, 0.0, 10)
	assert.Nil(t, ewmas)

	// empty list
	l = make([]float32, 0)
	ewmas = EwmaSeries32(l, 0.0, 10)
	assert.Nil(t, ewmas)

	// list < lookback
	l = []float32{
		1.0,
		9.0,
	}

	// should be = sma
	ewmas = EwmaSeries32(l, 0.0, 10)
	assert.Equal(t, float32(5.0), ewmas[len(ewmas)-1])

	// full data series on lb = 600, decay = default
	// read in test series file
	testseries := make([]float32, 0)
	f, err := os.Open("./mock/test_series.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		num, _ := strconv.ParseFloat(scanner.Text(), 32)
		testseries = append(testseries, float32(num))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ewmas = EwmaSeries32(testseries, 0.0, 600)
	assert.Equal(t, float32(39.855465), ewmas[len(ewmas)-1])

	// full data series on lb = 1200, decay = default
	ewmas = EwmaSeries32(testseries, 0.0, 1200)
	assert.Equal(t, float32(39.948578), ewmas[len(ewmas)-1])

}

func TestRollingEMA64(t *testing.T) {
	ewma := RollingEMA64(10, 5, 0.5)
	assert.Equal(t, 7.5, ewma)
}

func TestRollingEMA32(t *testing.T) {
	ewma := RollingEMA32(10, 5, 0.5)
	assert.Equal(t, float32(7.5), ewma)
}

// Benchmark tests
func BenchmarkSimpleAvg64(b *testing.B) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SimpleAvg64(xs)
	}
}

func BenchmarkStatsMean(b *testing.B) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stat.Mean(xs, nil)
	}
}

func BenchmarkVariance32(b *testing.B) {
	xs := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Variance32(xs)
	}
}

func BenchmarkStatsVariance(b *testing.B) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stat.Variance(xs, nil)
	}
}

func BenchmarkStdDev32(b *testing.B) {
	xs := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StdDev32(xs)
	}
}

func BenchmarkStatsStdDev(b *testing.B) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stat.StdDev(xs, nil)
	}
}

func BenchmarkEwmaSeries32(b *testing.B) {
	xs := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EwmaSeries32(xs, 0, 3)
	}
}
