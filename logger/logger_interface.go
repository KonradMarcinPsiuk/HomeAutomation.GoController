package logger

type LogOperator interface {
	Debug(msg string)
	Info(msg string)
	Trace(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Fatal(msg string)
	Panic(msg string)
	Close() error
}
