//go:build arm

package main

import (
	"GoController/gpio"
	"GoController/logger"
)

func initPinOperator(log logger.LogOperator) gpio.PinOperator {
	return gpio.NewRpiGPIO(log)
}
