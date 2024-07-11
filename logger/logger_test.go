package logger

import (
	"os"
	"strings"
	"testing"
)

func getConfig() LogConfig {
	return LogConfig{
		LogFilePath: "test.log",
		MaxSize:     5,
		MaxBackups:  2,
		MaxAge:      10,
		BufferSize:  1000,
		LogLevel:    DebugLevel,
	}
}

func TestLoggerWriteToFile(t *testing.T) {

	config := getConfig()

	if _, err := os.Stat(config.LogFilePath); err == nil {
		t.Fatalf("Log file %s already exists", config.LogFilePath)
	}

	t.Cleanup(func() {
		if err := os.Remove(config.LogFilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Cannot remove log file: %v", err)
		}
	})

	const (
		debugMessage = "debug message"
		infoMessage  = "info message"
		warnMessage  = "warn message"
		errorMessage = "error message"
	)

	logger := NewLogger(config)

	defer func(logger *ZeroLogLogger) {
		err := logger.Close()
		if err != nil {
			t.FailNow()
		}
	}(logger)

	logger.Debug(debugMessage)
	logger.Info(infoMessage)
	logger.Warn(warnMessage)
	logger.Error(errorMessage)

	err := logger.Close()
	if err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	if _, err := os.Stat(config.LogFilePath); os.IsNotExist(err) {
		t.Fatalf("Log file %s was not created", config.LogFilePath)
	}

	content, err := os.ReadFile(config.LogFilePath)
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
