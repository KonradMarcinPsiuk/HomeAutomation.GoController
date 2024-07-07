package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"os"
	"time"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

// NewLogger initializes a new ZerologLogger with configuration provided via LogConfig.
// It sets the global log level, configures log file rotation and compression,
// and combines console and file output into a multi-level writer.
func NewLogger(config LogConfig) *ZerologLogger {

	//Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	//Setup lumberjack to manage log files
	logFile := &lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	//Setup console writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	//Setup file writer
	diodeWriter := diode.NewWriter(logFile, config.BufferSize, config.FlushInterval, reportMissedLogs)

	//Setup one writer from console and file writers
	multi := zerolog.MultiLevelWriter(consoleWriter, diodeWriter)

	//Get log level and create logger instance
	logLevel := getZerologLevel(config.LogLevel)
	loggerInstance := zerolog.New(multi).With().Timestamp().Logger().Level(logLevel)

	return &ZerologLogger{logger: loggerInstance}
}

func reportMissedLogs(missed int) {
	fmt.Printf("Logger dropped %d messages\n", missed)
}

func getZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *ZerologLogger) logWithOptionalError(level zerolog.Level, msg string, errs ...error) {
	//Start a new message with the given level
	event := l.logger.WithLevel(level)

	//If error was passed to this function, write it to the log, otherwise just write the message
	if len(errs) > 0 && errs[0] != nil {
		event.Err(errs[0]).Msg(msg)
	} else {
		event.Msg(msg)
	}
}

func (l *ZerologLogger) Debug(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.DebugLevel, msg, errs...)
}

func (l *ZerologLogger) Info(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.InfoLevel, msg, errs...)
}

func (l *ZerologLogger) Warn(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.WarnLevel, msg, errs...)
}

func (l *ZerologLogger) Error(msg string, errs ...error) {
	l.logWithOptionalError(zerolog.ErrorLevel, msg, errs...)
}
