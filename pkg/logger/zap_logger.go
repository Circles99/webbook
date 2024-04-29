package logger

import "go.uber.org/zap"

// 适配器模式

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(args...)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(args...)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(args...)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(args...)...)
}

func (z *ZapLogger) toZapFields(args ...Field) []zap.Field {
	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		fields = append(fields, zap.Any(arg.Key, arg.Value))
	}
	return fields
}
