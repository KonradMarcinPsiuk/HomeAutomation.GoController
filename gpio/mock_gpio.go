//go:build !arm

package gpio

import (
	"GoController/logger"
	"fmt"
	"time"
)

type MockGPIO struct {
	gpioPin uint8
	name    string
	log     logger.LogOperator
}

func NewMockGPIO(logger logger.LogOperator) *MockGPIO {
	return &MockGPIO{name: "Mock GPIO Controller", log: logger}
}

func (g *MockGPIO) SetOutputPin(pin uint8) {
	g.log.Info(fmt.Sprintf("%s: Setting GPIO pin %v as output", g.name, pin))
	g.gpioPin = pin
}

func (g *MockGPIO) SetHigh() {
	g.log.Info(fmt.Sprintf("%s: GPIO Pin %v set to HIGH", g.name, g.gpioPin))
}

func (g *MockGPIO) SetLow() {
	g.log.Info(fmt.Sprintf("%s: GPIO Pin %v set to LOW", g.name, g.gpioPin))
}

func (g *MockGPIO) Open() error {
	g.log.Info(fmt.Sprintf("%s: Opening GPIO pin controller", g.name))

	time.Sleep(1 * time.Second)

	g.log.Info(fmt.Sprintf("%s: GPIO pin controller open", g.name))
	return nil
}

func (g *MockGPIO) Close() error {
	g.log.Info(fmt.Sprintf("%s: Closing GPIO pin controller", g.name))

	time.Sleep(1 * time.Second)

	g.log.Info(fmt.Sprintf("%s: GPIO pin controller closed", g.name))
	return nil
}
