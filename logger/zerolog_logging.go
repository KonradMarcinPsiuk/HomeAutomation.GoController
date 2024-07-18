// Package logger provides a logging utility with zerolog,
// including support for file rotation and in-memory buffering.
package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

// ZeroLogLogger wraps a zerolog.Logger with additional functionality
// such as file rotation and diode-based log buffering.
type ZeroLogLogger struct {
	logger           zerolog.Logger
	diodeWriter      diode.Writer
	lumberjackLogger *lumberjack.Logger
}

// NewLogger initializes a new ZeroLogLogger with configuration provided via LogConfig.
//
// It sets up a format for timestamping log entries and initializes log file rotation
// and buffering settings.
//
// Parameters:
//   - config: LogConfig struct containing logging configuration.
//
// Returns:
//   - *ZeroLogLogger: An instance of ZeroLogLogger.
func NewLogger(config LogConfig) *ZeroLogLogger {
	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Setup lumberjack to manage log files
	logFile := setupLumberjackLogger(&config)

	// Setup console writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}

	// Setup file writer
	diodeWriter := diode.NewWriter(logFile, config.BufferSize, config.PollInterval, reportMissedLogs)

	// Setup one writer from console and file writers
	multi := zerolog.MultiLevelWriter(consoleWriter, diodeWriter)

	// Create logger instance
	loggerInstance := zerolog.New(multi).With().Timestamp().Logger()

	return &ZeroLogLogger{logger: loggerInstance, lumberjackLogger: logFile, diodeWriter: diodeWriter}
}

// setupLumberjackLogger configures the lumberjack logger for log file rotation.
//
// Parameters:
//   - config: *LogConfig struct containing file rotation configuration.
//
// Returns:
//   - *lumberjack.Logger: An instance of lumberjack.Logger with the provided settings.
func setupLumberjackLogger(config *LogConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
}

// reportMissedLogs reports the number of dropped log messages due to buffer overflow.
//
// Parameters:
//   - missed: Number of missed log messages.
func reportMissedLogs(missed int) {
	fmt.Printf("Logger dropped %d messages\n", missed)
}

// Close flushes the diode writer and closes the lumberjack logger if it exists.
//
// If an error occurs during the closing process, it will be printed to standard output.
//
// Returns:
//   - error: Any error encountered during the closing process.
func (l *ZeroLogLogger) Close() error {
	// Flush the diode writer
	if err := l.diodeWriter.Close(); err != nil {
		return err
	}

	// Close the lumberjack logger
	if l.lumberjackLogger != nil {
		if err := l.lumberjackLogger.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Debug logs a debug-level message.
//
// Parameters:
//   - msg: The message to be logged.
func (l *ZeroLogLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Info logs an info-level message.
//
// Parameters:
//   - msg: The message to be logged.
func (l *ZeroLogLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Trace logs a trace-level message.
//
// Parameters:
//   - msg: The message to be logged.
func (l *ZeroLogLogger) Trace(msg string) {
	l.logger.Trace().Msg(msg)
}

// Warn logs a warning-level message.
//
// Parameters:
//   - msg: The message to be logged.
func (l *ZeroLogLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error logs an error-level message with an error.
//
// Parameters:
//   - msg: The message to be logged.
//   - err: The error associated with the message.
func (l *ZeroLogLogger) Error(msg string, err error) {
	l.logger.Error().Err(err).Msg(msg)
}

// Panic logs a panic-level message with an error.
//
// Parameters:
//   - msg: The message to be logged.
//   - err: The error associated with the message.
func (l *ZeroLogLogger) Panic(msg string, err error) {
	l.logger.Panic().Err(err).Msg(msg)
}
