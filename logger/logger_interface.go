package logger

type Logger interface {
	Debug(msg string, err ...error)
	Info(msg string, err ...error)
	Warn(msg string, err ...error)
	Error(msg string, err ...error)
}
