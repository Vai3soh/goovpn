package logger

import (
	"fmt"
	"os"

	"github.com/fangdingjun/go-log/v5"
)

type InterfaceLogger interface {
	log.Interface
	log.StdLog
}

type Logger struct {
	*log.FixedSizeFileWriter
	*log.Logger
}

type OptionLogger func(*Logger)

func NewLogger(opts ...OptionLogger) *Logger {
	f := &Logger{Logger: new(log.Logger)}

	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithLogWriteFile(logfile *string, logFileCount *int, logFileSize *int64) OptionLogger {
	return func(l *Logger) {
		l.Name = *logfile
		l.MaxCount = *logFileCount
		l.MaxSize = *logFileSize * 1024 * 1024
		log.Default.Out = l.FixedSizeFileWriter
		l.Out = log.Default.Out
	}
}

func WithLogLevel(loglevel *string) OptionLogger {
	return func(l *Logger) {
		ld, err := log.ParseLevel(*loglevel)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		log.Default.Level = ld
		l.Level = log.Default.Level
		if l.Out == nil {
			l.Out = log.Default.Out
		}
	}
}

func WithLogTextFormatter() OptionLogger {
	return func(l *Logger) {
		l.Formatter = new(log.TextFormatter)
	}
}

func WithLogJsonFormatter() OptionLogger {
	return func(l *Logger) {
		l.Formatter = new(log.JSONFormatter)
	}
}
