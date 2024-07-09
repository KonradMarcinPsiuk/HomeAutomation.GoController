package mqtt

import (
	"GoController/logger"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// MQTTConfig represents the configuration for an MQTT client. It contains the following fields:
//
// - Broker: The address of the MQTT broker to connect to.
// - Topic: The topic to subscribe to or publish messages to.
// - ClientID: The ID to use when connecting to the MQTT broker.
// - Username: The username to use for authentication.
// - Password: The password to use for authentication.
type MQTTConfig struct {
	Broker   string
	Topic    string
	ClientID string
	Username string
	Password string
}

// PublishMessage Message to be published
type PublishMessage struct {
	Topic   string
	Payload []byte
	QoS     byte
}

// ReceivedMessage represents a message received from an MQTT broker.
// It contains the topic and payload of the received message.
type ReceivedMessage struct {
	Topic   string
	Payload []byte
}

// MQTTClient is a type that represents an MQTT client.
// Note: MQTTClient requires the MQTTConfig and logger.LogOperator types to be properly initialized.
type MQTTClient struct {
	client         mqtt.Client
	config         MQTTConfig
	logger         logger.LogOperator
	msgQueue       chan PublishMessage
	messageHandler mqtt.MessageHandler
}

// NewMQTTClient initializes a new MQTTClient
func NewMQTTClient(config MQTTConfig, logOperator logger.LogOperator) *MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.KeepAlive = 10
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.ConnectRetryInterval = 5 * time.Second

	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		logOperator.Error("Connection lost")
	}
	opts.OnConnect = func(mqtt.Client) {
		logOperator.Info("Connection established")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		logOperator.Info("Attempting to reconnect")
	}

	client := mqtt.NewClient(opts)
	mc := &MQTTClient{
		client:   client,
		config:   config,
		msgQueue: make(chan PublishMessage, 100),
		logger:   logOperator,
	}

	go mc.connect()
	go mc.startProcessing()

	return mc
}

// Connect attempts to connect the MQTT client to the broker
func (c *MQTTClient) connect() {

	c.logger.Info(fmt.Sprintf("Connecting to broker: %s", c.config.Broker))

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		c.logger.Error(fmt.Sprintf("Error connecting to broker: %v", token.Error()))
		return
	}

	c.logger.Info("Connected to broker")

	c.logger.Info(fmt.Sprintf("Subscribing to topic: %s", c.config.Topic))

	// Subscribe to the topic upon successful connection
	if token := c.client.Subscribe(c.config.Topic, 2, c.messageHandler); token.Wait() && token.Error() != nil {
		c.logger.Error(fmt.Sprintf("Error subscribing to topic: %v", token.Error()))
	}

	c.logger.Info(fmt.Sprintf("Subscribed to topic: %s", c.config.Topic))
}

// Publish sends an MQTT message to the specified topic with the given payload and quality of service (QoS) level.
// The message is added to the message queue and will be processed asynchronously.
func (c *MQTTClient) Publish(topic string, payload []byte, qos byte) {
	c.msgQueue <- PublishMessage{Topic: topic, Payload: payload, QoS: qos}
}

// startProcessing processes messages from the queue and publishes them
func (c *MQTTClient) startProcessing() {
	c.logger.Info("Starting message processing")

	for msg := range c.msgQueue {
		if c.client.IsConnected() {
			token := c.client.Publish(msg.Topic, msg.QoS, false, msg.Payload)
			token.Wait()
		} else {
			c.logger.Error("Client not connected, could not publish message")
		}
	}
}

// SetMessageCallback sets the callback function to handle incoming MQTT messages.
// The callback function should have the signature func(message ReceivedMessage),
// where ReceivedMessage is a struct containing the topic and payload of the received message.
// When a message is received, it is converted to ReceivedMessage and passed to the callback function.
func (c *MQTTClient) SetMessageCallback(callback func(message ReceivedMessage)) {
	c.messageHandler = func(client mqtt.Client, msg mqtt.Message) {
		// Convert the mqtt.Message to ReceivedMessage
		receivedMsg := ReceivedMessage{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}

		callback(receivedMsg)
	}
}
