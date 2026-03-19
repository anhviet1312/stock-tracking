package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel int) (*zap.SugaredLogger, error) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder,
			EncodeLevel:  CustomLevelEncoder,
			EncodeTime:   SyslogTimeEncoder,
		},
	}

	logger, _ := cfg.Build()

	// nolint:errcheck
	defer logger.Sync() // flushes buffer, if any

	return logger.Sugar(), nil
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
