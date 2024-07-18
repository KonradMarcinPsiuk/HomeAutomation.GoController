package logger

import "time"

type LogConfig struct {
	LogFilePath  string
	BufferSize   int
	PollInterval time.Duration
	MaxSize      int
	MaxBackups   int
	MaxAge       int
}
