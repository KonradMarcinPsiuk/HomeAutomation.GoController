package mqtt

type MessageReceiver interface {
	SetMessageCallback(func(message ReceivedMessage))
}
