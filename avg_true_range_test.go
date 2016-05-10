package technical

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrueRange64(t *testing.T) {
	s := []float64{0, 0, 0, 0}
	assert.Equal(t, 0.0, TrueRange64(s, 10.0))

	s = []float64{0, 0, 0, 0}
	assert.Equal(t, 0.0, TrueRange64(s, 0.0))

	s = []float64{1, 0, 4, 2, 7, 9, 4}
	assert.Equal(t, 8.0, TrueRange64(s, 0.0))
	assert.Equal(t, 10.0, TrueRange64(s, 11.0))
}

func TestTrueRange32(t *testing.T) {
	s := []float32{0, 0, 0, 0}
	assert.Equal(t, float32(0.0), TrueRange32(s, 10.0))

	s = []float32{0, 0, 0, 0}
	assert.Equal(t, float32(0.0), TrueRange32(s, 0.0))

	s = []float32{1, 4, 0, 2, 7, 9, 4}
	assert.Equal(t, float32(8.0), TrueRange32(s, float32(0.0)))
	assert.Equal(t, float32(10.0), TrueRange32(s, float32(11.0)))
}

func TestRollingATR64(t *testing.T) {
	assert.Equal(t, 0.0, RollingATR64(8.0, 3.0, 0))
	assert.Equal(t, 7.0, RollingATR64(8.0, 3.0, 5))
}

func TestRollingATR32(t *testing.T) {
	assert.Equal(t, float32(0.0), RollingATR32(float32(8.0), float32(3.0), 0))
	assert.Equal(t, float32(7.0), RollingATR32(float32(8.0), float32(3.0), 5))
}

func TestStaticATR64(t *testing.T) {
	s := []float64{1, 4, 2, 7, 9, 4}

	// 0 n value
	assert.Equal(t, 0.0, StaticATR64(s, 0, 10))

	// 0 s value
	assert.Equal(t, 0.0, StaticATR64(s, 10, 0))

	// not even data, should be 0
	assert.Equal(t, 0.0, StaticATR64(s, 10, 10))

	// just the true range of the first period
	assert.Equal(t, 8.0, StaticATR64(s, 3, 5))

	s = []float64{1, 4, 2, 7, 9, 4, 5, 6, 7, 9, 4, 7, 6, 9, 2, 3, 10, 12, 11, 15}
	assert.Equal(t, 8.777777777777779, StaticATR64(s, 3, 5))

	// full day series test
	// read in test series file
	testseries := make([]float64, 0)
	f, err := os.Open("./mock/test_atr_series.txt")
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
	assert.Equal(t, 0.25532787677004043, StaticATR64(testseries, 30, 300))
}

func TestStaticATR32(t *testing.T) {
	s := []float32{1, 4, 2, 7, 9, 4}

	// 0 n value
	assert.Equal(t, float32(0.0), StaticATR32(s, 0, 10))

	// 0 s value
	assert.Equal(t, float32(0.0), StaticATR32(s, 10, 0))

	// not even data, should be 0
	assert.Equal(t, float32(0.0), StaticATR32(s, 10, 10))

	// just the true range of the first period
	assert.Equal(t, float32(8.0), StaticATR32(s, 3, 5))

	s = []float32{1, 4, 2, 7, 9, 4, 5, 6, 7, 9, 4, 7, 6, 9, 2, 3, 10, 12, 11, 15}
	assert.Equal(t, float32(8.777777777777779), StaticATR32(s, 3, 5))

	// full day series test
	// read in test series file
	testseries := make([]float32, 0)
	f, err := os.Open("./mock/test_atr_series.txt")
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

	assert.Equal(t, float32(0.25532743), StaticATR32(testseries, 30, 300))
}
