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
func NewLogger(config LogConfig) *ZeroLogLogger {

	//Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	//Setup lumberjack to manage log files
	logFile := setupLumberjackLogger(&config)

	//Setup console writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}

	//Setup file writer
	diodeWriter := diode.NewWriter(logFile, config.BufferSize, config.PollInterval, reportMissedLogs)

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

// Close flushes the diode writer and closes the lumberjack logger if it exists. If an error occurs during the closing process,
// it will be printed to standard output.
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

// Debug logs a debug-level message.
func (l *ZeroLogLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Info logs an info-level message.
func (l *ZeroLogLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Trace logs a trace-level message.
func (l *ZeroLogLogger) Trace(msg string) {
	l.logger.Trace().Msg(msg)
}

// Warn logs a warning-level message.
func (l *ZeroLogLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error logs an error-level message with error.
func (l *ZeroLogLogger) Error(msg string, err error) {
	l.logger.Error().Err(err).Msg(msg)
}

// Panic logs a panic-level message with error.
func (l *ZeroLogLogger) Panic(msg string) {
	l.logger.Panic().Msg(msg)
}

// Fatal logs a fatal-level message with error.
func (l *ZeroLogLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}
