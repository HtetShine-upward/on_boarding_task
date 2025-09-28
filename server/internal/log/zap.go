package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	// JSON 既定、UTCタイムスタンプなど Production の標準を踏襲
	cfg.Level = zap.NewAtomicLevelAt(parseLevel(level))
	lg, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return lg
}

func parseLevel(s string) zapcore.Level {
	switch strings.ToLower(s) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}

}
