//go:build !arm

package main

import "GoController/gpio"

func initPinOperator() gpio.PinOperator {
	return gpio.NewMockGPIO()
}
