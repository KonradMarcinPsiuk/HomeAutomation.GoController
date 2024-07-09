package main

import (
	"GoController/logger"
	"GoController/mqtt"
	"fmt"
	"os"
	"sync"
	"time"
)

const logFileName string = "logs/goController.log"

func main() {

	var mutex sync.Mutex

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

	mqttConfig := mqtt.MQTTConfig{
		Broker:   os.Getenv("mqtt_broker"),
		Topic:    "topic",
		ClientID: "go-mqtt-client",
		Username: "",
		Password: "",
	}

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

	messageCallback := func(message mqtt.ReceivedMessage) {

		mutex.Lock()
		defer mutex.Unlock()
		log.Info(fmt.Sprint("Received message: ", string(message.Payload)))

		if string(message.Payload) == "SetHigh" {
			pinOperator.SetHigh()
		}
		if string(message.Payload) == "SetLow" {
			pinOperator.SetLow()
		}
	}

	mqttClient := mqtt.NewMQTTClient(mqttConfig, log)

	mqttClient.SetMessageCallback(messageCallback)

	select {}
}
