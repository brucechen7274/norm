package logger

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := New(os.Stdout, Config{
		Colorful: true,
		LogLevel: DebugLevel,
	})
	ctx := context.Background()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, "error message")
	logger.Trace(ctx, &TraceRecord{
		NGQL: `GO FROM "player102" OVER serve YIELD dst(edge);`,
		Err:  errors.New("error"),
	})

	logger = logger.LogMode(WarnLevel)
	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Trace(ctx, &TraceRecord{
		NGQL: `GO FROM "player102" OVER serve YIELD dst(edge);`,
	})

	logger = logger.LogMode(SilentLevel)
	logger.Error(ctx, "error message")
}
