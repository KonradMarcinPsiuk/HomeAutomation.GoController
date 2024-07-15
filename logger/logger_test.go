package logger

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// getConfig returns a LogConfig struct with the following default values:
func getConfig() LogConfig {
	return LogConfig{
		LogFilePath: "test.log",
		MaxSize:     5,
		MaxBackups:  2,
		MaxAge:      10,
		BufferSize:  1000,
	}
}

// TestLoggerWriteToFile is a test function that verifies the functionality of writing log messages to a file.
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
		panicErrMessage = "panic_err_message"
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

	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered in TestLoggerWriteToFile", r)
				done <- true
			}
		}()
		logger.Panic("panic message", errors.New("panic_err_message"))
	}()
	<-done

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

	checkLogMessage(t, logContent, traceMessage, "", "trace")
	checkLogMessage(t, logContent, debugMessage, "", "debug")
	checkLogMessage(t, logContent, infoMessage, "", "info")
	checkLogMessage(t, logContent, warnMessage, "", "warn")
	checkLogMessage(t, logContent, errorMessage, errorErrMessage, "error")
	checkLogMessage(t, logContent, panicMessage, panicErrMessage, "panic")
}

func checkLogMessage(t *testing.T, logContent, message, errMessage, logType string) {
	if !strings.Contains(logContent, message) || (errMessage != "" && !strings.Contains(logContent, errMessage)) {
		t.Errorf("Log file does not contain the %s message: %s", logType, message)
	} else {
		t.Logf("Log file contains the %s message", logType)
	}
}
