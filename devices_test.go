package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDevices(t *testing.T) {
	devices := getDevices()

	for _, d := range devices {
		fmt.Printf("%s: \n", d.Path)
		fmt.Printf("Name:\t%s\n", d.Name)
		fmt.Printf("Brightness:\t%d\n", d.Brightness())
		fmt.Printf("Max Brightness:\t%d\n", d.MaxBrightness())
		fmt.Printf("Percentage:\t%d\n", d.BrightnessPercentage())

		if d.Name == "intel_backlight" {
			d.SetBrightnessPercentage(80)
		}
	}
}

func TestDevice_SetBrightnessPercentage(t *testing.T) {
	d := getDevices()
	for _, d := range d {
		if d.Name == "intel_backlight" {
			d.SetBrightnessPercentage(50)
		}
	}
}

func TestDevice_SetSmoothBrightnessPercentage(t *testing.T) {
	d := getDevices()
	for _, d := range d {
		if d.Name == "intel_backlight" {
			d.SetSmoothBrightnessPercentage(60, 100*time.Millisecond)
			assert.Equal(t, uint(60), d.BrightnessPercentage())
			d.SetSmoothBrightnessPercentage(40, 100*time.Millisecond)
			assert.Equal(t, uint(40), d.BrightnessPercentage())
			d.SetSmoothBrightnessPercentage(60, 100*time.Millisecond)
			assert.Equal(t, uint(60), d.BrightnessPercentage())
		}
	}
}
