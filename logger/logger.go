package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

type Level int

const (
	DebugLevel Level = iota // By default, trace information of query execution is recorded as debug logs.
	InfoLevel
	WarnLevel
	ErrorLevel
	SilentLevel // Below this level, norm will not print any logs or trace information.
)

type TraceRecord struct {
	NGQL string
	Err  error
}

type Interface interface {
	// LogMode returns a logger with the specified log level.
	LogMode(level Level) Interface

	// Trace is used to track the execution of Nebula Graph queries.
	// In the default logger, trace records are printed as debug logs.
	Trace(ctx context.Context, record *TraceRecord)

	Debug(ctx context.Context, msg string, data ...any)
	Info(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
}

type defaultLogger struct {
	Config
	stdLog                             *log.Logger
	debugStr, infoStr, warnStr, errStr string
}

type Config struct {
	Colorful bool
	LogLevel Level
}

var Default = New(os.Stdout, Config{
	Colorful: true,
	LogLevel: WarnLevel,
})

func New(writer io.Writer, conf Config) Interface {
	logger := &defaultLogger{
		Config:   conf,
		stdLog:   log.New(writer, "\r\n", log.Ldate|log.Ltime|log.Lshortfile),
		debugStr: "[DEBUG] ",
		infoStr:  "[INFO] ",
		warnStr:  "[WARN] ",
		errStr:   "[ERROR] ",
	}
	if conf.Colorful {
		logger.debugStr = "\x1b[36m[DEBUG] \x1b[0m"
		logger.infoStr = "\x1b[32m[INFO] \x1b[0m"
		logger.warnStr = "\x1b[33m[WARN] \x1b[0m"
		logger.errStr = "\x1b[31m[ERROR] \x1b[0m"
	}
	return logger
}

func (l *defaultLogger) LogMode(level Level) Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, data ...any) {
	l.message(ctx, DebugLevel, l.debugStr+msg, data...)
}

func (l *defaultLogger) Info(ctx context.Context, msg string, data ...any) {
	l.message(ctx, InfoLevel, l.infoStr+msg, data...)
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.message(ctx, WarnLevel, l.warnStr+msg, data...)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, data ...any) {
	l.message(ctx, ErrorLevel, l.errStr+msg, data...)
}

func (l *defaultLogger) message(_ context.Context, level Level, msg string, data ...any) {
	if level < l.LogLevel {
		return
	}
	_ = l.stdLog.Output(3, fmt.Sprintf(msg, data...))
}

func (l *defaultLogger) Trace(ctx context.Context, record *TraceRecord) {
	if l.LogLevel >= SilentLevel || record == nil {
		return
	}
	if record.Err == nil {
		l.message(ctx, DebugLevel, l.debugStr+"[norm] nGQL: %s", record.NGQL)
	} else {
		l.message(ctx, DebugLevel, l.debugStr+"[norm] nGQL: %s err: %v", record.NGQL, record.Err)
	}
}
