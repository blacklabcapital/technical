package technical

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundBoundToNearestCent64(t *testing.T) {
	b := Bound64{10.123, 0.0, 10.123}

	RoundBoundToNearestCent64(&b)

	// upper
	assert.Equal(t, 10.13, b.Upper)

	// lower
	assert.Equal(t, 10.12, b.Lower)
}

func TestRoundBoundToNearestCent32(t *testing.T) {
	b := Bound32{float32(10.123), float32(0.0), float32(10.123)}

	RoundBoundToNearestCent32(&b)

	// upper
	assert.Equal(t, float32(10.13), b.Upper)

	// lower
	assert.Equal(t, float32(10.12), b.Lower)
}

func TestCompareBound64(t *testing.T) {
	firstBound := &Bound64{
		Lower:    -5,
		Midpoint: 0,
		Upper:    5,
	}
	secondBound := &Bound64{
		Lower:    2,
		Midpoint: 4,
		Upper:    6,
	}

	res := CompareBound64(firstBound, secondBound)
	assert.True(t, res)

	res = CompareBound64(secondBound, firstBound)
	assert.False(t, res)
}

func TestCompareBound32(t *testing.T) {
	firstBound := &Bound32{
		Lower:    -5,
		Midpoint: 0,
		Upper:    5,
	}
	secondBound := &Bound32{
		Lower:    2,
		Midpoint: 4,
		Upper:    6,
	}

	res := CompareBound32(firstBound, secondBound)
	assert.True(t, res)

	res = CompareBound32(secondBound, firstBound)
	assert.False(t, res)
}

func TestBollBound64(t *testing.T) {
	var b Bound64

	// nil list
	var l []float64
	b = BollBound64(l, 0.0, 10)
	assert.Equal(t, Bound64{}, b)

	// empty list
	l = make([]float64, 0)
	b = BollBound64(l, 0.0, 10)
	assert.Equal(t, Bound64{}, b)

	l = []float64{
		1,
		3,
		5,
		7,
		9,
	}
	b = BollBound64(l, 4, 2)
	assert.Equal(t, 4.0, b.Midpoint)
	assert.Equal(t, -1.6568542494923806, b.Lower)
	assert.Equal(t, 9.65685424949238, b.Upper)
}

func TestBollBound32(t *testing.T) {
	var b Bound32

	// nil list
	var l []float32
	b = BollBound32(l, 0.0, 10)
	assert.Equal(t, Bound32{}, b)

	// empty list
	l = make([]float32, 0)
	b = BollBound32(l, 0.0, 10)
	assert.Equal(t, Bound32{}, b)

	l = []float32{
		1,
		3,
		5,
		7,
		9,
	}
	b = BollBound32(l, 4, 2)
	assert.Equal(t, float32(4.0), b.Midpoint)
	assert.Equal(t, float32(-1.6568542), b.Lower)
	assert.Equal(t, float32(9.656855), b.Upper)
}

func TestRollingBollingerConst64(t *testing.T) {
	var b Bound64

	l := []float64{
		1,
		3,
		5,
		7,
		9,
	}

	// static k
	b = RollingBollingerConst64(l, 4.0, 2)
	assert.Equal(t, 4.0, b.Midpoint)
	assert.Equal(t, -1.6568542494923806, b.Lower)
	assert.Equal(t, 9.65685424949238, b.Upper)

}

func TestRollingBollingerConst32(t *testing.T) {
	var b Bound32

	l := []float32{
		1,
		3,
		5,
		7,
		9,
	}

	// static k
	b = RollingBollingerConst32(l, 4.0, 2)
	assert.Equal(t, float32(4.0), b.Midpoint)
	assert.Equal(t, float32(-1.6568542), b.Lower)
	assert.Equal(t, float32(9.656855), b.Upper)

}

func TestRollingBollingerSMA64(t *testing.T) {
	var b Bound64

	l := []float64{
		1,
		3,
		5,
		7,
		9,
	}

	// simple average
	b = RollingBollingerSMA64(l, 2)
	assert.Equal(t, 5.0, b.Midpoint)
	assert.Equal(t, -0.6568542494923806, b.Lower)
	assert.Equal(t, 10.65685424949238, b.Upper)

}

func TestRollingBollingerSMA32(t *testing.T) {
	var b Bound32

	l := []float32{
		1,
		3,
		5,
		7,
		9,
	}

	// simple average
	b = RollingBollingerSMA32(l, 2)
	assert.Equal(t, float32(5.0), b.Midpoint)
	assert.Equal(t, float32(-0.65685415), b.Lower)
	assert.Equal(t, float32(10.656855), b.Upper)

}

func TestRollingBollingerEMA64(t *testing.T) {
	var b Bound64

	l := []float64{
		1,
		3,
		5,
		7,
		9,
	}

	l2 := []float64{
		3,
		5,
		7,
		9,
		11,
	}

	l3 := []float64{
		5,
		7,
		9,
		11,
		13,
	}

	// simple average
	b = RollingBollingerSMA64(l, 2)
	assert.Equal(t, 5.0, b.Midpoint)
	assert.Equal(t, -0.6568542494923806, b.Lower)
	assert.Equal(t, 10.65685424949238, b.Upper)

	// rolling ewma 1
	b = RollingBollingerEMA64(l2, 1/3.0, 13.0, b.Midpoint, 2)
	assert.Equal(t, 7.666666666666667, b.Midpoint)
	assert.Equal(t, 2.0098124171742864, b.Lower)
	assert.Equal(t, 13.323520916159048, b.Upper)

	// rolling ewma 2
	b = RollingBollingerEMA64(l3, 1/3.0, 15.0, b.Midpoint, 2)
	assert.Equal(t, 10.11111111111111, b.Midpoint)
	assert.Equal(t, 4.45425686161873, b.Lower)
	assert.Equal(t, 15.76796536060349, b.Upper)

	// test full test series on rolling ewma
	// read in test series file
	var testseries []float64
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

	lb := 1200               // lookback
	y := 2.0 / float64(lb+1) // lambda
	for i, v := range testseries {
		if i < lb {
			continue
		}

		period := testseries[i-lb : i]
		if i == lb { // first bound is simple avg
			b = RollingBollingerSMA64(period, 2.0)
		} else {
			b = RollingBollingerEMA64(period, y, v, b.Midpoint, 2.0)
		}
	}

	assert.Equal(t, 39.9488124468039, b.Midpoint)
}

func TestRollingBollingerEMA32(t *testing.T) {
	var b Bound32

	l := []float32{
		1,
		3,
		5,
		7,
		9,
	}

	l2 := []float32{
		3,
		5,
		7,
		9,
		11,
	}

	l3 := []float32{
		5,
		7,
		9,
		11,
		13,
	}

	// simple average
	b = RollingBollingerSMA32(l, 2)
	assert.Equal(t, float32(5.0), b.Midpoint)
	assert.Equal(t, float32(-0.65685415), b.Lower)
	assert.Equal(t, float32(10.656855), b.Upper)

	// rolling ewma 1
	b = RollingBollingerEMA32(l2, 1/float32(3), 13.0, b.Midpoint, 2)
	assert.Equal(t, float32(7.6666665), b.Midpoint)
	assert.Equal(t, float32(2.0098124), b.Lower)
	assert.Equal(t, float32(13.323521), b.Upper)

	// rolling ewma 2
	b = RollingBollingerEMA32(l3, 1/float32(3), 15.0, b.Midpoint, 2)
	assert.Equal(t, float32(10.111111), b.Midpoint)
	assert.Equal(t, float32(4.4542565), b.Lower)
	assert.Equal(t, float32(15.767965), b.Upper)

	// read in test series file
	var testseries []float32
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
	// test full test series on rolling ewma
	lb := 1200               // lookback
	y := 2.0 / float32(lb+1) // lambda
	for i, v := range testseries {
		if i < lb {
			continue
		}

		period := testseries[i-lb : i]
		if i == lb { // first bound is simple avg
			b = RollingBollingerSMA32(period, 2.0)
		} else {
			b = RollingBollingerEMA32(period, y, v, b.Midpoint, 2.0)
		}
	}

	assert.Equal(t, float32(39.948578), b.Midpoint)

}

func TestStaticBollingerConst64(t *testing.T) {
	var band []Bound64

	l := []float64{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band = StaticBollingerConst64(l, 5, 0.0, 2.0)
	assert.Equal(t, 0.0, band[5].Midpoint)
	assert.Equal(t, -5.656854249492381, band[5].Lower)
	assert.Equal(t, 5.656854249492381, band[5].Upper)

}

func TestStaticBollingerConst32(t *testing.T) {
	var band []Bound32

	l := []float32{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band = StaticBollingerConst32(l, 5, 0.0, 2.0)
	assert.Equal(t, float32(0.0), band[5].Midpoint)
	assert.Equal(t, float32(-5.656854249), band[5].Lower)
	assert.Equal(t, float32(5.656854249), band[5].Upper)

}

func TestStaticBollingerSMA64(t *testing.T) {
	var band []Bound64

	l := []float64{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band = StaticBollingerSMA64(l, 5, 2)
	assert.Equal(t, 7.0, band[5].Midpoint)
	assert.Equal(t, 1.3431457505076194, band[5].Lower)
	assert.Equal(t, 12.65685424949238, band[5].Upper)
}

func TestStaticBollingerSMA32(t *testing.T) {
	var band []Bound32

	l := []float32{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band = StaticBollingerSMA32(l, 5, 2)
	assert.Equal(t, float32(7), band[5].Midpoint)
	assert.Equal(t, float32(1.34314585), band[5].Lower)
	assert.Equal(t, float32(12.65685425), band[5].Upper)
}

func TestStaticBollingerEMA64(t *testing.T) {
	var band []Bound64

	// read in test series file
	var testseries []float64
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

	// test full test series
	lb := 1200 // lookback
	y := 0.0   // use default smoothing
	band = StaticBollingerEMA64(testseries, lb, y, 2)
	assert.Equal(t, 39.9488124468039, band[len(band)-1].Midpoint)
}

func TestStaticBollingerEMA32(t *testing.T) {
	var band []Bound32

	// read in test series file
	var testseries []float32
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

	// test full test series
	lb := 1200        // lookback
	y := float32(0.0) // use default smoothing
	band = StaticBollingerEMA32(testseries, lb, y, 2)
	assert.Equal(t, float32(39.948578), band[len(band)-1].Midpoint)
}
