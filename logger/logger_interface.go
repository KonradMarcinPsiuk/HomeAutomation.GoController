package logger

type LogOperator interface {
	Debug(msg string, err ...error)
	Info(msg string, err ...error)
	Warn(msg string, err ...error)
	Error(msg string, err ...error)
}

type Logger interface {
	NewLogger(config LogConfig) LogOperator
}
