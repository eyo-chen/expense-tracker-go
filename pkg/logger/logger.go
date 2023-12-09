package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	*zap.Logger
}

// registration logger
var regLogger Logger

// Register registers a new logger
func Register() {
	regLogger = newZapLogger()
}

func newZapLogger() Logger {
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.DebugLevel)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
	)

	return &ZapLogger{zap.New(core)}
}

func (l *ZapLogger) Debugw(msg string, options ...interface{}) {
	l.Sugar().Debugw(msg, options...)
}

func (l *ZapLogger) Infow(msg string, options ...interface{}) {
	l.Sugar().Infow(msg, options...)
}

func (l *ZapLogger) Warnw(msg string, options ...interface{}) {
	l.Sugar().Warnw(msg, options...)
}

func (l *ZapLogger) Errorw(msg string, options ...interface{}) {
	l.Sugar().Errorw(msg, options...)
}

func (l *ZapLogger) Panicw(msg string, options ...interface{}) {
	l.Sugar().Panicw(msg, options...)
}

func (l *ZapLogger) Fatalw(msg string, options ...interface{}) {
	l.Sugar().Fatalw(msg, options...)
}
