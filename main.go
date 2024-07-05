package main

import (
	"GoController/logger"
	"fmt"
	"time"
)

const logFileName string = "logs/goController.log"

func main() {

	config := logger.Config{
		LogFilePath:   logFileName,
		BufferSize:    1000,
		FlushInterval: 1 * time.Second,
		MaxSize:       1,
		MaxBackups:    3,
		MaxAge:        30,
	}

	log := logger.NewZerologLogger(config)

	log.Error("Testing log", fmt.Errorf("a new error occurred at %v", time.Now()))

	pinOperator := initPinOperator(log)

	pinOperatorOpenErr := pinOperator.Open()
	if pinOperatorOpenErr != nil {
		log.Error("Failed to open pin operator")
	}

	defer func() {
		pinOperatorCloseError := pinOperator.Close()
		if pinOperatorCloseError != nil {
			log.Error("Failed to close pin operator")
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
