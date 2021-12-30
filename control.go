package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type IncreaseCmd struct {
	To       string `arg:"positional"`
	SmoothMs uint   `arg:"-s,--smooth-ms" default:"0"`
}

type DecreaseCmd struct {
	To       string `arg:"positional"`
	SmoothMs uint   `arg:"-s,--smooth-ms" default:"0"`
}

type SetCmd struct {
	To       string `arg:"positional"`
	SmoothMs uint   `arg:"-s,--smooth-ms" default:"0"`
}

type ListCmd struct{}

func doIncrease(cmd *Args) error {
	if cmd.Increase == nil {
		return fmt.Errorf("increase cannot be nil")
	}

	d := selectDevice(cmd)
	if d == nil {
		return fmt.Errorf("device not found")
	}

	if cmd.Increase.To == "" {
		// Increase by 5%
		cmd.Increase.To = "5%"
	}

	return execIncDec(d, cmd.Increase.To, true, cmd.Increase.SmoothMs)
}

func doSet(cmd *Args) error {
	if cmd.Set == nil {
		return fmt.Errorf("increase cannot be nil")
	}

	d := selectDevice(cmd)
	if d == nil {
		return fmt.Errorf("device not found")
	}

	if cmd.Set.To == "" {
		// Increase by 5%
		cmd.Set.To = "5%"
	}

	return execSet(d, cmd.Set.To, cmd.Set.SmoothMs)
}

func execSet(d *Device, to string, smoothMs uint) error {
	v, err := strconv.ParseInt(
		strings.Replace(to, "%", "", -1),
		10,
		32,
	)
	if err != nil {
		return fmt.Errorf("unable to parse value")
	}
	if v < 0 {
		return fmt.Errorf("invalid value %s, must be positive", to)
	}

	if strings.HasSuffix(to, "%") {
		// Percentage
		targetPercentage := correctPercentage(int(v))

		if smoothMs == 0 {
			d.SetBrightnessPercentage(targetPercentage)
		} else {
			d.SetSmoothBrightnessPercentage(targetPercentage,
				time.Duration(smoothMs)*time.Millisecond,
			)
		}

	} else {
		// Absolute value
		targetValue := correctValue(int(v), d.MaxBrightness())

		if smoothMs == 0 {
			d.SetBrightness(targetValue)
		} else {
			d.SetSmoothBrightness(targetValue, time.Duration(smoothMs)*time.Millisecond)
		}
	}
	return nil
}

func correctValue(value int, maxValue uint) uint {
	if value > int(maxValue) {
		return maxValue
	}

	if value < 0 {
		return 0
	}

	return uint(value)
}

func correctPercentage(percentage int) uint {
	return correctValue(percentage, 100)
}

func execIncDec(d *Device, to string, plus bool, smoothMs uint) error {
	v, err := strconv.ParseInt(
		strings.Replace(to, "%", "", -1),
		10,
		32,
	)
	if err != nil {
		return fmt.Errorf("unable to parse value")
	}
	if v < 0 {
		return fmt.Errorf("invalid value %s, must be positive", to)
	}

	if strings.HasSuffix(to, "%") {
		// Percentage
		targetPercentage := int(d.BrightnessPercentage())
		if plus {
			targetPercentage += int(v)
		} else {
			targetPercentage -= int(v)
		}

		thePercentage := correctPercentage(targetPercentage)

		if smoothMs == 0 {
			d.SetBrightnessPercentage(thePercentage)
		} else {
			d.SetSmoothBrightnessPercentage(thePercentage,
				time.Duration(smoothMs)*time.Millisecond,
			)
		}

	} else {
		// Absolute value
		targetValue := int(d.Brightness())
		if plus {
			targetValue += int(v)
		} else {
			targetValue -= int(v)
		}

		theValue := correctValue(targetValue, d.MaxBrightness())

		if smoothMs == 0 {
			d.SetBrightness(theValue)
		} else {
			d.SetSmoothBrightness(theValue, time.Duration(smoothMs)*time.Millisecond)
		}
	}
	return nil
}

func doList(cmd *Args) error {
	if cmd.List == nil {
		return fmt.Errorf("list cannot be nil")
	}

	for _, d := range getDevices() {
		fmt.Printf("%s\n", d.Name)
	}
	return nil
}

func selectDevice(cmd *Args) *Device {
	devices := getDevices()
	var d *Device = nil
	if len(devices) == 0 {
		return nil
	}

	if cmd.Device == "" {
		d = &devices[0]
	}
	for _, c := range devices {
		if c.Name == cmd.Device {
			d = &c
			break
		}
	}

	return d
}

func doDecrease(cmd *Args) error {
	if cmd.Decrease == nil {
		return fmt.Errorf("decrease cannot be nil")
	}
	d := selectDevice(cmd)
	if d == nil {
		return fmt.Errorf("device not found")
	}

	if cmd.Decrease.To == "" {
		// Increase by 5%
		cmd.Decrease.To = "5%"
	}

	return execIncDec(d, cmd.Decrease.To, false, cmd.Decrease.SmoothMs)
}
