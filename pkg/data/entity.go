package data

import (
	"math"

	"github.com/google/uuid"
)

// randomEntityId randomly generates a new entity ID string
// in the form of a UUID.
func randomEntityId() string {
	return uuid.New().String()
}

// roundFloat64 rounds a float64 to the dp precision.
func roundFloat64(n float64, dp uint) float64 {
	mod := math.Pow(10, float64(dp))
	return math.Round(n*mod) / mod
}
