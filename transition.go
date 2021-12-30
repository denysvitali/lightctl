package main

import (
	"math"
	"time"
)

func getStepTime(from uint, to uint, duration time.Duration) time.Duration {
	// 10 - 11 - 12 - 13 - 14 - 15
	// |_________________________|
	//           500 ms

	// 10 - 15 = 5 steps
	// 500 ms / 5 steps = 100ms/step

	distance := math.Abs(float64(int(to) - int(from)))
	durationFloat := math.Round(float64(duration.Milliseconds()) / distance)
	return time.Duration(durationFloat) * time.Millisecond
}
