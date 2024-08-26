package mqtt

type MessageReceiver interface {
	SetMessageCallback(func(message ReceivedMessage))
	Publish(topic string, payload []byte, qos byte)
}
