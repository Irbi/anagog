package tools

import (
	"math"
	"math/rand"
)

func RandGeo(min, max float64) float64 {
	res := toFixed(min+rand.Float64()*(max-min), 5)
	return res
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
