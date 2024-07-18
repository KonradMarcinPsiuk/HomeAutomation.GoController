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
		LogFilePath:  logFileName,
		BufferSize:   1000,
		PollInterval: 1 * time.Second,
		MaxSize:      1,
		MaxBackups:   3,
		MaxAge:       30,
	}

	log := logger.NewLogger(logConfig)

	mqttConfig := mqtt.MQTTConfig{
		Broker:   os.Getenv("mqtt_broker"),
		Topic:    "topic",
		ClientID: "go-mqtt-client",
		Username: "",
		Password: "",
	}

	log.Info("Starting Pin Operator")
	pinOperator := initPinOperator(log)
	err := pinOperator.Open()

	if err != nil {
		log.Error("Failed to open pin operator", err)
		if err := pinOperator.Open(); err != nil {
			log.Panic("Failed to open pin operator", err)
		}

		defer func() {
			err = pinOperator.Close()
			if err != nil {
				log.Error("Failed to close pin operator", err)
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
}
