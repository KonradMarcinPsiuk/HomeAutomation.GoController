package logger

import "time"

type Config struct {
	LogFilePath   string
	BufferSize    int
	FlushInterval time.Duration
	MaxSize       int
	MaxBackups    int
	MaxAge        int
}
