package logger

import (
	"errors"
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

	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in TestLoggerWriteToFile", r)
		}
	}()

	const (
		traceMessage    = "trace message"
		debugMessage    = "debug message"
		infoMessage     = "info message"
		warnMessage     = "warn message"
		errorMessage    = "error message"
		errorErrMessage = "err_err_message"
		panicMessage    = "panic message"
		fatalMessage    = "fatal message"
	)

	logger := NewLogger(config)

	defer func(logger *ZeroLogLogger) {
		err := logger.Close()
		if err != nil {
			t.FailNow()
		}
	}(logger)

	logger.Trace(traceMessage)
	logger.Debug(debugMessage)
	logger.Info(infoMessage)
	logger.Warn(warnMessage)
	logger.Error(errorMessage, errors.New(errorErrMessage))
	logger.Panic(panicMessage)
	logger.Fatal(fatalMessage)

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
	if !strings.Contains(logContent, traceMessage) {
		t.Errorf("Log file does not contain the trace message: %s", traceMessage)
	}
	if !strings.Contains(logContent, debugMessage) {
		t.Errorf("Log file does not contain the debug message: %s", debugMessage)
	}
	if !strings.Contains(logContent, infoMessage) {
		t.Errorf("Log file does not contain the info message: %s", infoMessage)
	}
	if !strings.Contains(logContent, warnMessage) {
		t.Errorf("Log file does not contain the warn message: %s", warnMessage)
	}
	if !strings.Contains(logContent, errorMessage) || !strings.Contains(logContent, errorErrMessage) {
		t.Errorf("Log file does not contain the error message: %s", errorMessage)
	}
	if !strings.Contains(logContent, panicMessage) {
		t.Errorf("Log file does not contain the panic message: %s", panicMessage)
	}
	if !strings.Contains(logContent, fatalMessage) {
		t.Errorf("Log file does not contain the fatal message: %s", fatalMessage)
	}
}
