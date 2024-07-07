package mqtt

import (
	"GoController/logger"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Config struct {
	Broker   string
	Topic    string
	ClientID string
	Username string
	Password string
}

type MQTTClient struct {
	client mqtt.Client
	config Config
	logger logger.Logger
}

func NewMQTTClient(config Config) *MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	return &MQTTClient{client: client, config: config}
}

func (c *MQTTClient) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to broker: %v", token.Error())
	}
	return nil
}

func (c *MQTTClient) Subscribe() error {
	if token := c.client.Subscribe(c.config.Topic, 2, messageSubHandler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error subscribing to topic: %v", token.Error())
	}
	return nil
}

func (c *MQTTClient) Publish() {
	go func() {
		for {
			text := "Hello MQTT QoS 2"
			token := c.client.Publish(c.config.Topic, 2, false, text)
			token.Wait()
			time.Sleep(3 * time.Second)
		}
	}()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
