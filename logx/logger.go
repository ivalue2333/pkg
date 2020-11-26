package logx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var gLogger Logger

func SetGLogger(l1 Logger) {
	gLogger = l1
}

func GetGLogger() Logger {
	return gLogger
}

func NewLogger(opts ...Option) (l Logger, err error) {
	return NewLoggerWithOptions(newOptions(opts...))
}

func NewLoggerWithOptions(options Options) (l Logger, err error) {
	l = NewDefaultLogger()
	if err = initLoggerWithOptions(l, options); err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}
	return l, nil
}

func initLoggerWithOptions(l Logger, options Options) (err error) {
	// 如果配置里指定了日志等级，则解析并设置，否则默认等级是info。
	if options.Level != "" {
		level, err := ParseLevel(options.Level)
		if err != nil {
			return errors.Wrapf(err, "failed to parse level(%s)", options.Level)
		}
		l.SetLevel(level)
	}
	return
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	gLogger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	gLogger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	gLogger.Fatalf(ctx, format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	gLogger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	gLogger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	gLogger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	gLogger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	gLogger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	gLogger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	gLogger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	gLogger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	gLogger.Fatal(ctx, args...)
}

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}
