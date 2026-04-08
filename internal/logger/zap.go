package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.SugaredLogger
}

func (z *ZapLogger) Info(msg string, args ...any) {
	z.l.Infow(msg, args...)
}

func (z *ZapLogger) Error(msg string, args ...any) {
	z.l.Errorw(msg, args...)
}

func (z *ZapLogger) Debug(msg string, args ...any) {
	z.l.Debugw(msg, args...)
}

func (z *ZapLogger) Warn(msg string, args ...any) {
	z.l.Warnw(msg, args...)
}

func (z *ZapLogger) Close() error {
	return z.l.Sync()
}

func NewZapLogger(l *zap.Logger) ILogger {
	return &ZapLogger{
		l: l.Sugar(),
	}
}
