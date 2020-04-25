package utils

import (
	"math/rand"
	"time"
)

// RandFloats .
func RandFloats(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	r := min + rand.Float64()*(max-min)
	return r
}
