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

	logger.Debug(ctx, "hello world")
	logger.Info(ctx, "hello world")
	logger.Warn(ctx, "hello world")
	logger.Error(ctx, "hello world")

	logger.Trace(ctx, &TraceRecord{
		NGQL: `GO FROM "player102" OVER serve YIELD dst(edge);`,
		Err:  errors.New("error"),
	})
}
