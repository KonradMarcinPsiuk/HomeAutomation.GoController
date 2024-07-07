package main

import (
	"GoController/logger"
	"GoController/mqtt"
	"time"
)

const logFileName string = "logs/goController.log"

func main() {

	logConfig := logger.LogConfig{
		LogFilePath:   logFileName,
		BufferSize:    1000,
		FlushInterval: 1 * time.Second,
		MaxSize:       1,
		MaxBackups:    3,
		MaxAge:        30,
		LogLevel:      logger.InfoLevel,
	}

	log := logger.NewLogger(logConfig)

	mqttConfig := mqtt.Config{
		Broker:   "tcp://192.168.1.110:1883",
		Topic:    "your/topic",
		ClientID: "go-mqtt-client",
		Username: "",
		Password: "",
	}
	mqttClient := mqtt.NewMQTTClient(mqttConfig)

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
