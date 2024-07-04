//go:build arm

package gpio

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
)

type RpiGPIO struct {
	gpioPin rpio.Pin
	name    string
}

func NewRpiGPIO() *RpiGPIO {
	return &RpiGPIO{name: "Raspberry Pi GPIO Controller"}
}

func (g *RpiGPIO) SetOutputPin(pin uint8) {
	g.gpioPin = rpio.Pin(pin)
	g.gpioPin.Output()
}

func (g *RpiGPIO) SetHigh() {
	g.gpioPin.High()
	fmt.Printf("%s: GPIO Pin %v set to HIGH\n", g.name, g.gpioPin)
}

func (g *RpiGPIO) SetLow() {
	g.gpioPin.Low()
	fmt.Printf("%s: GPIO Pin %v set to LOW\n", g.name, g.gpioPin)
}

func (g *RpiGPIO) Open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	return nil
}

func (g *RpiGPIO) Close() error {
	err := rpio.Close()
	if err != nil {
		return err
	}
	return nil
}
