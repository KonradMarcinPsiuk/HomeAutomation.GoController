package gpio

type PinOperator interface {
	SetOutputPin(pin uint8)
	SetHigh()
	SetLow()
	Open() error
	Close() error
}
