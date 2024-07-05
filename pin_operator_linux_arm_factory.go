//go:build arm

package main

import (
	"GoController/gpio"
	"GoController/logger"
)

func initPinOperator(logger logger.Logger) gpio.PinOperator {
	return gpio.NewRpiGPIO(logger)
}
