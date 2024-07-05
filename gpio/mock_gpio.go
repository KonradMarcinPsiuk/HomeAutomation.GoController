//go:build !arm

package gpio

import (
	"GoController/logger"
	"fmt"
)

type MockGPIO struct {
	gpioPin uint8
	name    string
	logger  logger.Logger
}

func NewMockGPIO(logger logger.Logger) *MockGPIO {
	return &MockGPIO{name: "Mock GPIO Controller", logger: logger}
}

func (g *MockGPIO) SetOutputPin(pin uint8) {
	g.logger.Info(fmt.Sprintf("%s: Setting GPIO pin %v as output", g.name, pin))
	g.gpioPin = pin
}

func (g *MockGPIO) SetHigh() {
	g.logger.Info(fmt.Sprintf("%s: GPIO Pin %v set to HIGH", g.name, g.gpioPin))
}

func (g *MockGPIO) SetLow() {
	g.logger.Info(fmt.Sprintf("%s: GPIO Pin %v set to LOW", g.name, g.gpioPin))
}

func (g *MockGPIO) Open() error {
	g.logger.Info(fmt.Sprintf("%s: Opening GPIO pin controller", g.name))

	return nil
}

func (g *MockGPIO) Close() error {
	g.logger.Info(fmt.Sprintf("%s: Closing GPIO pin controller", g.name))

	return nil
}
