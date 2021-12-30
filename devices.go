package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Device struct {
	Name          string
	Path          string
	Subsystem     string
	maxBrightness uint
	canOpen       bool
}

var devicesClasses = []string{"backlight", "leds"}

func getDevices() []Device {
	var deviceList []Device
	basePath := "/sys/class/"
	for _, class := range devicesClasses {
		// For each directory, get the devices
		currentDir := path.Join(basePath, class)
		f, err := ioutil.ReadDir(currentDir)
		if err != nil {
			// skip
			logrus.Warnf("unable to list contents in directory %s: %v", currentDir, err)
			continue
		}

		for _, e := range f {
			deviceType := e.Mode().Type()
			if deviceType == fs.ModeSymlink {
				thePath := path.Join(currentDir, e.Name())
				deviceList = append(deviceList, Device{
					Name:      e.Name(),
					Path:      thePath,
					Subsystem: class,
					canOpen:   true,
				})
			}
		}
	}

	return deviceList
}

func (d *Device) Brightness() uint {
	brightnessPath := path.Join(d.Path, "brightness")
	brightness, err := os.ReadFile(brightnessPath)
	if err != nil {
		return 0
	}

	brightnessInt, err := strconv.ParseInt(strings.TrimSpace(string(brightness)), 10, 32)
	if err != nil {
		return 0
	}

	return uint(brightnessInt)
}

func (d *Device) MaxBrightness() uint {
	if d.maxBrightness != 0 {
		return d.maxBrightness
	}
	brightnessPath := path.Join(d.Path, "max_brightness")
	brightness, err := os.ReadFile(brightnessPath)
	if err != nil {
		return 0
	}

	brightnessInt, err := strconv.ParseInt(strings.TrimSpace(string(brightness)), 10, 32)
	if err != nil {
		return 0
	}

	d.maxBrightness = uint(brightnessInt)
	return d.maxBrightness
}

func (d *Device) BrightnessPercentage() uint {
	return uint(math.Round(100 * float64(d.Brightness()) / float64(d.MaxBrightness())))
}

func (d *Device) SetBrightnessPercentage(p uint) {
	value := uint(math.Round(float64(p) * float64(d.MaxBrightness()) / 100.0))

	d.writeBrightness(value)
}

func (d *Device) writeBrightness(value uint) {

	// Try to write brightness directly
	var err error
	if d.canOpen {
		err := d.writeBrightnessSys(value)
		if err == nil {
			return
		}
		d.canOpen = false
		logrus.Debugf("using DBUS from now on for device %s", d.Name)
	}

	err = d.writeBrightnessDbus(value)
	if err != nil {
		logrus.Fatalf("unable to set brightness via DBUS: %v", err)
	}
}

func (d *Device) writeBrightnessSys(value uint) error {
	brightnessFilePath := path.Join(d.Path, "brightness")
	f, err := os.OpenFile(brightnessFilePath, os.O_WRONLY, 0)
	if err != nil {
		logrus.Debugf("unable to open brightness path for writing: %v", err)
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("%d", value))
	if err != nil {
		return err
	}
	return f.Close()
}

func (d *Device) writeBrightnessDbus(value uint) error {
	conn, err := dbus.ConnectSystemBus()
	defer conn.Close()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	obj := conn.Object("org.freedesktop.login1", "/org/freedesktop/login1/session/auto")
	err = obj.Call("org.freedesktop.login1.Session.SetBrightness", 0, d.Subsystem, d.Name, uint32(value)).Store()
	logrus.Debugf("Set brightness via DBUS: (%v, %v, %v)",
		d.Subsystem,
		d.Name,
		uint32(value),
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) SetSmoothBrightnessPercentage(target uint, duration time.Duration) {
	bp := d.BrightnessPercentage()
	if bp == target {
		return
	}

	bm := d.MaxBrightness()

	targetValue := uint(math.Round(float64(target) / 100 * float64(bm)))
	d.SetSmoothBrightness(targetValue, duration)
}

func (d *Device) SetBrightness(value uint) {
	d.writeBrightness(value)
}

func (d *Device) SetSmoothBrightness(target uint, duration time.Duration) {
	b := d.Brightness()
	if b == target {
		return
	}

	stepTime := getStepTime(b, target, duration)
	distance := getDistance(b, target)
	stepSize := 1

	if stepTime < 5*time.Millisecond {
		stepTime = 5 * time.Millisecond
		maxSteps := int(duration / stepTime)
		stepSize = int(math.Floor(float64(distance) / float64(maxSteps)))
	}

	if target > b {
		// Increasing
		for i := b; i <= target; i += uint(stepSize) {
			d.SetBrightness(i)
			time.Sleep(stepTime)
		}
	} else {
		// Decreasing
		for i := b; i >= target; i -= uint(stepSize) {
			d.SetBrightness(i)
			time.Sleep(stepTime)
		}
	}
	d.SetBrightness(target)
}

func getDistance(from uint, to uint) uint {
	return uint(math.Abs(float64(from) - float64(to)))
}
