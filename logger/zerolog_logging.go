package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type ZeroLogLogger struct {
	logger           zerolog.Logger
	diodeWriter      diode.Writer
	lumberjackLogger *lumberjack.Logger
}

// NewLogger initializes a new ZeroLogLogger with configuration provided via LogConfig.
// It sets the global log level, configures log file rotation and compression,
// and combines console and file output into a multi-level writer.
func NewLogger(config LogConfig) *ZeroLogLogger {

	//Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	//Setup lumberjack to manage log files
	logFile := setupLumberjackLogger(&config)

	//Setup console writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}

	//Setup file writer
	diodeWriter := diode.NewWriter(logFile, config.BufferSize, config.FlushInterval, reportMissedLogs)

	//Setup one writer from console and file writers
	multi := zerolog.MultiLevelWriter(consoleWriter, diodeWriter)

	//Create logger instance
	loggerInstance := zerolog.New(multi).With().Timestamp().Logger()

	return &ZeroLogLogger{logger: loggerInstance, lumberjackLogger: logFile, diodeWriter: diodeWriter}
}

func setupLumberjackLogger(config *LogConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
}

func reportMissedLogs(missed int) {
	fmt.Printf("Logger dropped %d messages\n", missed)
}

func (l *ZeroLogLogger) logWithOptionalError(level zerolog.Level, msg string, errs ...error) {
	//Start a new message with the given level
	event := l.logger.WithLevel(level)

	//If error was passed to this function, write it to the log, otherwise just write the message
	if len(errs) > 0 && errs[0] != nil {
		event.Err(errs[0]).Msg(msg)
	} else {
		event.Msg(msg)
	}
}

// Close flushes the diode writer and closes the lumberjack logger if it exists. If an error occurs during the closing process,
// it will be printed to standard output.
func (l *ZeroLogLogger) Close() error {
	// Flush the diode writer
	var err = l.diodeWriter.Close()
	if err != nil {
		return err
	}

	// Close the lumberjack logger
	if l.lumberjackLogger != nil {
		lumberjackErr := l.lumberjackLogger.Close()
		if lumberjackErr != nil {
			return lumberjackErr
		}
	}

	return nil
}

// Debug logs a debug-level message with optional errors.
func (l *ZeroLogLogger) Debug(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.DebugLevel, msg, errs...)
}

// Info logs an info-level message with optional errors.
func (l *ZeroLogLogger) Info(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.InfoLevel, msg, errs...)
}

// Warn logs a warning-level message with optional errors.
func (l *ZeroLogLogger) Warn(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.WarnLevel, msg, errs...)
}

// Error logs an error-level message with optional errors.
func (l *ZeroLogLogger) Error(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.ErrorLevel, msg, errs...)
}
