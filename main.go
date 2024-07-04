package main

import (
	"time"
)

func main() {

	pinOperator := initPinOperator()

	pinOperatorErr := pinOperator.Open()

	if pinOperatorErr != nil {
		panic(pinOperatorErr)
	}

	pinOperator.SetOutputPin(10)

	for i := 0; i < 10; i++ {

		time.Sleep(1 * time.Second)
		pinOperator.SetHigh()

		time.Sleep(1 * time.Second)
		pinOperator.SetLow()

	}

	pinOperatorCloseError := pinOperator.Close()

	if pinOperatorCloseError != nil {
		panic(pinOperatorCloseError)
	}

}
