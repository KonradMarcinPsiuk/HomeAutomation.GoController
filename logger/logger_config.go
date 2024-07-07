package logger

import "time"

type LogConfig struct {
	LogFilePath   string
	BufferSize    int
	FlushInterval time.Duration
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	LogLevel      LogLevel
}
