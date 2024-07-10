package logger

import (
	"os"
	"strings"
	"testing"
)

func TestLoggerWriteToFile(t *testing.T) {
	const logFilePath = "test.log"

	const debugMessage = "debug message"
	const infoMessage = "info message"
	const warnMessage = "warn message"
	const errorMessage = "error message"

	config := LogConfig{
		LogFilePath: logFilePath,
		BufferSize:  1000,
		MaxSize:     1, // 1 MB
		MaxBackups:  1,
		MaxAge:      1,
		LogLevel:    DebugLevel,
	}

	logger := NewLogger(config)

	logger.Debug(debugMessage)
	logger.Info(infoMessage)
	logger.Warn(warnMessage)
	logger.Error(errorMessage)

	logger.Close()

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, debugMessage) || !strings.Contains(logContent, infoMessage) ||
		!strings.Contains(logContent, warnMessage) || !strings.Contains(logContent, errorMessage) {
		t.Errorf("Log file does not contain the expected content")
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Cannot remove log file: %v", err)
		}
	}(logFilePath)
}
