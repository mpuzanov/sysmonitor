package calc

import (
	"math"
	"math/rand"
	"time"
)

// RandFloats случайное значение из диапазона
func RandFloats(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	r := min + rand.Float64()*(max-min)
	return r
}

// RoundToFixed округление с заданным значением точности
func RoundToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}
