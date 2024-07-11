package logger

import (
	"os"
	"strings"
	"testing"
)

func TestLoggerWriteToFile(t *testing.T) {
	const logFilePath = "test.log"

	const (
		debugMessage = "debug message"
		infoMessage  = "info message"
		warnMessage  = "warn message"
		errorMessage = "error message"
	)

	config := LogConfig{
		LogFilePath: logFilePath,
		BufferSize:  1000,
		MaxSize:     1,
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

	t.Cleanup(func() {
		err := os.Remove(logFilePath)
		if err != nil {
			t.Fatalf("Cannot remove log file: %v", err)
		}
	})
}
