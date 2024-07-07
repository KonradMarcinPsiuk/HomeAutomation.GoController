//go:build arm

package gpio

import (
	"GoController/logger"
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
)

type RpiGPIO struct {
	gpioPin rpio.Pin
	name    string
	logger  logger.LogOperator
}

func NewRpiGPIO(logger logger.LogOperator) *RpiGPIO {
	return &RpiGPIO{name: "Raspberry Pi GPIO Controller", logger: logger}
}

func (g *RpiGPIO) SetOutputPin(pin uint8) {
	g.gpioPin = rpio.Pin(pin)
	g.gpioPin.Output()
}

func (g *RpiGPIO) SetHigh() {
	g.gpioPin.High()
	g.logger.Info(fmt.Sprintf("%s: GPIO Pin %v set to HIGH", g.name, g.gpioPin))
}

func (g *RpiGPIO) SetLow() {
	g.gpioPin.Low()
	g.logger.Info(fmt.Sprintf("%s: GPIO Pin %v set to LOW", g.name, g.gpioPin))
}

func (g *RpiGPIO) Open() error {
	g.logger.Info(fmt.Sprintf("%s: Opening GPIO pin controller", g.name))

	err := rpio.Open()

	if err != nil {
		return err
	}

	g.logger.Error(fmt.Sprintf("%s: GPIO pin controller open", g.name), err)
	return nil
}

func (g *RpiGPIO) Close() error {
	g.logger.Info(fmt.Sprintf("%s: Closing GPIO pin controller", g.name))

	err := rpio.Close()

	if err != nil {
		return err
	}

	g.logger.Info(fmt.Sprintf("%s: GPIO pin controller closed", g.name))
	return nil
}
