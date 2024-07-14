package gpio

import (
	"GoController/logger"
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
)

type RpiGPIO struct {
	gpioPin rpio.Pin
	name    string
	log     logger.LogOperator
}

func NewRpiGPIO(log logger.LogOperator) *RpiGPIO {
	return &RpiGPIO{name: "Raspberry Pi GPIO Controller", log: log}
}

func (g *RpiGPIO) SetOutputPin(pin uint8) {
	g.log.Info(fmt.Sprintf("%s: Setting GPIO pin %v as output", g.name, pin))

	g.gpioPin = rpio.Pin(pin)
	g.gpioPin.Output()

	g.log.Info(fmt.Sprintf("%s: GPIO pin %v set as output", g.name, g.gpioPin))
}

func (g *RpiGPIO) Open() error {
	g.log.Info(fmt.Sprintf("%s: Opening GPIO pin controller", g.name))

	err := rpio.Open()

	if err != nil {
		return err
	}

	g.log.Info(fmt.Sprintf("%s: GPIO pin controller open", g.name))
	return nil
}

func (g *RpiGPIO) SetHigh() {
	g.gpioPin.High()
	g.log.Info(fmt.Sprintf("%s: GPIO Pin %v set to HIGH", g.name, g.gpioPin))
}

func (g *RpiGPIO) SetLow() {
	g.gpioPin.Low()
	g.log.Info(fmt.Sprintf("%s: GPIO Pin %v set to LOW", g.name, g.gpioPin))
}

func (g *RpiGPIO) Close() error {
	g.log.Info(fmt.Sprintf("%s: Closing GPIO pin controller", g.name))

	err := rpio.Close()

	if err != nil {
		return err
	}

	g.log.Info(fmt.Sprintf("%s: GPIO pin controller closed", g.name))
	return nil
}
