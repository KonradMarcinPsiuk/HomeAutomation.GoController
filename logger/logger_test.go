package logger

import (
	"os"
	"strings"
	"testing"
)

func TestLoggerWriteToFile(t *testing.T) {
	const logFilePath = "test.log"

	if _, err := os.Stat(logFilePath); err == nil {
		t.Fatalf("Log file %s already exists", logFilePath)
	}

	t.Cleanup(func() {
		if err := os.Remove(logFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Cannot remove log file: %v", err)
		}
	})

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

	defer logger.Close()

	logger.Debug(debugMessage)
	logger.Info(infoMessage)
	logger.Warn(warnMessage)
	logger.Error(errorMessage)

	logger.Close()

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatalf("Log file %s was not created", logFilePath)
	}

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, debugMessage) {
		t.Errorf("Log file does not contain the debug message: %s", debugMessage)
	}
	if !strings.Contains(logContent, infoMessage) {
		t.Errorf("Log file does not contain the info message: %s", infoMessage)
	}
	if !strings.Contains(logContent, warnMessage) {
		t.Errorf("Log file does not contain the warn message: %s", warnMessage)
	}
	if !strings.Contains(logContent, errorMessage) {
		t.Errorf("Log file does not contain the error message: %s", errorMessage)
	}
}
