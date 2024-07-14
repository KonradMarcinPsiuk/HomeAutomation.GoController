package main

import (
	"GoController/logger"
	"fmt"
	"time"
)

const logFileName string = "logs/goController.log"

func main() {

	logConfig := logger.LogConfig{
		LogFilePath:  logFileName,
		BufferSize:   1000,
		PollInterval: 1 * time.Second,
		MaxSize:      1,
		MaxBackups:   3,
		MaxAge:       30,
	}

	log := logger.NewLogger(logConfig)

	log.Info("Starting Pin Operator")
	pinOperator := initPinOperator(log)
	if err := pinOperator.Open(); err != nil {
		log.Panic("Failed to open pin operator", err)
	}

	defer func() {
		if err := pinOperator.Close(); err != nil {
			log.Panic("Failed to close pin operator", err)
		}

		if err := log.Close(); err != nil {
			fmt.Printf("Failed to close logger. Error: %s\n", err)
		}
	}()

	pinOperator.SetOutputPin(10)

	for {
		time.Sleep(1 * time.Second)
		pinOperator.SetHigh()

		time.Sleep(1 * time.Second)
		pinOperator.SetLow()
	}
}
