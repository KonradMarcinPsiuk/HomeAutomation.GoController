package main

import (
	"log"
	"time"
)

func main() {

	pinOperator := initPinOperator()

	pinOperatorErr := pinOperator.Open()
	if pinOperatorErr != nil {
		log.Fatalf("Failed to open pin operator: %v", pinOperatorErr)
	}

	defer func() {
		pinOperatorCloseError := pinOperator.Close()
		if pinOperatorCloseError != nil {
			log.Fatalf("Failed to close pin operator: %v", pinOperatorCloseError)
		}
	}()

	pinOperator.SetOutputPin(10)

	for i := 0; i < 10; i++ {

		time.Sleep(1 * time.Second)
		pinOperator.SetHigh()

		time.Sleep(1 * time.Second)
		pinOperator.SetLow()

	}
}
