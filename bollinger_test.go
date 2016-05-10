package bollinger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleAvg32(t *testing.T) {
	var (
		v   float32
		err error
	)

	// nil list
	var l []float32
	_, err = SimpleAvg32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// empty list
	l = make([]float32, 0)
	_, err = SimpleAvg32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	l = append(l, 2.0)
	l = append(l, 4.0)

	// good
	v, err = SimpleAvg32(l)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, float32(3.0), v)
}

func TestVariance32(t *testing.T) {
	var (
		err error
		v   float32
	)

	// nil list
	var l []float32
	_, err = Variance32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// empty list
	l = make([]float32, 0)
	_, err = Variance32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// good
	l = []float32{
		1,
		3,
		5,
		7,
	}

	v, err = Variance32(l)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(5.0), v)
}

func TestStdDev32(t *testing.T) {
	var (
		err error
		v   float32
	)

	// nil list
	var l []float32
	_, err = StdDev32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// empty list
	l = make([]float32, 0)
	_, err = StdDev32(l)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// good
	l = []float32{
		10,
		100,
	}

	v, err = StdDev32(l)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(45), v)
}

func TestCompareBound32(t *testing.T) {
	firstBound := &Bound32{
		Lower:    -5,
		MidPoint: 0,
		Upper:    5,
	}
	secondBound := &Bound32{
		Lower:    2,
		MidPoint: 4,
		Upper:    6,
	}

	res := CompareBound32(firstBound, secondBound)
	assert.True(t, res)

	res = CompareBound32(secondBound, firstBound)
	assert.False(t, res)
}

func TestEwmaSeries32(t *testing.T) {
	var (
		err   error
		ewmas []float32
	)

	// nil list
	var l []float32
	_, err = EwmaSeries32(l, 0.0, 10)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// empty list
	l = make([]float32, 0)
	_, err = EwmaSeries32(l, 0.0, 10)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// list < lookback
	l = []float32{
		1.0,
		9.0,
	}

	_, err = EwmaSeries32(l, 0.0, 10)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// full data series on lb = 600, decay = default
	ewmas, err = EwmaSeries32(TESTSERIESA, 0.0, 600)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(39.855465), ewmas[len(ewmas)-1])

	// full data series on lb = 1200, decay = default
	ewmas, err = EwmaSeries32(TESTSERIESA, 0.0, 1200)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(39.948578), ewmas[len(ewmas)-1])

}

func TestRollingEMA32(t *testing.T) {
	ewma := RollingEMA32(10, 5, 0.5)
	assert.Equal(t, float32(7.5), ewma)
}

func TestBollBound32(t *testing.T) {
	var (
		err error
		b   Bound32
	)

	// nil list
	var l []float32
	_, err = BollBound32(l, 0.0, 10)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	// empty list
	l = make([]float32, 0)
	_, err = BollBound32(l, 0.0, 10)
	assert.NotNil(t, err)
	if err == nil {
		return
	}

	l = []float32{
		1,
		3,
		5,
		7,
		9,
	}
	b, err = BollBound32(l, 4, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(4.0), b.MidPoint)
	assert.Equal(t, float32(-1.6568542), b.Lower)
	assert.Equal(t, float32(9.656855), b.Upper)
}

func TestRollingBollingerConst32(t *testing.T) {
	var (
		err error
		b   Bound32
	)

	l := []float32{
		1,
		3,
		5,
		7,
		9,
	}

	// static k
	b, err = RollingBollingerConst32(l, 4.0, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(4.0), b.MidPoint)
	assert.Equal(t, float32(-1.6568542), b.Lower)
	assert.Equal(t, float32(9.656855), b.Upper)

}

func TestRollingBollingerSMA32(t *testing.T) {
	var (
		err error
		b   Bound32
	)

	l := []float32{
		1,
		3,
		5,
		7,
		9,
	}

	// simple average
	b, err = RollingBollingerSMA32(l, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(5.0), b.MidPoint)
	assert.Equal(t, float32(-0.65685415), b.Lower)
	assert.Equal(t, float32(10.656855), b.Upper)

}

func TestRollingBollingerEMA32(t *testing.T) {
	var (
		err error
		b   Bound32
	)

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
	b, err = RollingBollingerSMA32(l, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(5.0), b.MidPoint)
	assert.Equal(t, float32(-0.65685415), b.Lower)
	assert.Equal(t, float32(10.656855), b.Upper)

	// rolling ewma 1
	b, err = RollingBollingerEMA32(l2, 1/float32(3), 13.0, b.MidPoint, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(7.6666665), b.MidPoint)
	assert.Equal(t, float32(2.0098124), b.Lower)
	assert.Equal(t, float32(13.323521), b.Upper)

	// rolling ewma 2
	b, err = RollingBollingerEMA32(l3, 1/float32(3), 15.0, b.MidPoint, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(10.111111), b.MidPoint)
	assert.Equal(t, float32(4.4542565), b.Lower)
	assert.Equal(t, float32(15.767965), b.Upper)

	// test full test series on rolling ewma
	lb := 1200               // lookback
	y := 2.0 / float32(lb+1) // lambda
	for i, v := range TESTSERIESA {
		if i < lb {
			continue
		}

		period := TESTSERIESA[i-lb : i]
		if i == lb { // first bound is simple avg
			b, err = RollingBollingerSMA32(period, 2.0)
			if err != nil {
				return
			}
		} else {
			b, err = RollingBollingerEMA32(period, y, v, b.MidPoint, 2.0)
			if err != nil {
				return
			}
		}
	}

	assert.Equal(t, float32(39.948578), b.MidPoint)

}

func TestStaticBollingerConst32(t *testing.T) {
	var (
		err  error
		band []Bound32
	)

	l := []float32{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band, err = StaticBollingerConst32(l, 5, 0.0, 2.0)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, float32(0.0), band[5].MidPoint)
	assert.Equal(t, float32(-5.656854249), band[5].Lower)
	assert.Equal(t, float32(5.656854249), band[5].Upper)

}

func TestStaticBollingerSMA32(t *testing.T) {
	var (
		err  error
		band []Bound32
	)

	l := []float32{
		1,
		3,
		5,
		7,
		9,
		11,
		13,
	}

	band, err = StaticBollingerSMA32(l, 5, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, float32(5), band[5].MidPoint)
	assert.Equal(t, float32(-0.65685415), band[5].Lower)
	assert.Equal(t, float32(10.65685425), band[5].Upper)
}

func TestStaticBollingerEMA32(t *testing.T) {
	var (
		err  error
		band []Bound32
	)

	// test full test series
	lb := 1200        // lookback
	y := float32(0.0) // use default smoothing
	band, err = StaticBollingerEMA32(TESTSERIESA, lb, y, 2)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, float32(39.948578), band[len(band)-1].MidPoint)
}
