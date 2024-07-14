package logger

type LogOperator interface {
	Debug(msg string)
	Info(msg string)
	Trace(msg string)
	Warn(msg string)
	Error(msg string, err ...error)
	Fatal(msg string, err ...error)
	Panic(msg string, err ...error)
	Close() error
}
