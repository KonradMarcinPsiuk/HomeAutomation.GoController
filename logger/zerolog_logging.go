package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger(config Config) *ZerologLogger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logFile := &lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	diodeWriter := diode.NewWriter(logFile, config.BufferSize, config.FlushInterval,
		func(missed int) {
			log.Error().Msgf("Logger dropped %d messages", missed)
		})

	multi := zerolog.MultiLevelWriter(consoleWriter, diodeWriter)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return &ZerologLogger{logger: log.Logger}
}

func (l *ZerologLogger) logWithOptionalError(level zerolog.Level, msg string, errs ...error) {
	event := l.logger.WithLevel(level)

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
