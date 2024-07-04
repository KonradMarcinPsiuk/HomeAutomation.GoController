//go:build !arm

package gpio

import (
	"fmt"
)

type MockGPIO struct {
	gpioPin uint8
	name    string
}

func NewMockGPIO() *MockGPIO {
	return &MockGPIO{name: "Mock GPIO Controller"}
}

func (g *MockGPIO) SetOutputPin(pin uint8) {
	g.gpioPin = pin
}

func (g *MockGPIO) SetHigh() {
	fmt.Printf("%s: GPIO Pin %v set to HIGH\n", g.name, g.gpioPin)
}

func (g *MockGPIO) SetLow() {
	fmt.Printf("%s: GPIO Pin %v set to LOW\n", g.name, g.gpioPin)
}

func (g *MockGPIO) Open() error {
	fmt.Println("Opening GPIO pin controller")

	return nil
}

func (g *MockGPIO) Close() error {
	fmt.Println("Closing GPIO pin controller")

	return nil
}
