package logx

import (
	"context"
	"github.com/pkg/errors"
)

var logger Logger

func SetLogger(l1 Logger) {
	logger = l1
}

func GetLogger() Logger {
	return logger
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
	// 如果配置里指定了日志文件，则解析并设置，否则默认写到stderr。
	if options.File != "" {
		err = HandleFileOutput(l, options.File)
		if err != nil {
			errors.Wrapf(err, "failed to set logger.Output and set flow_control")
		}
	}
	return
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	logger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	logger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	logger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logger.Fatalf(ctx, format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	logger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	logger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	logger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	logger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	logger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	logger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	logger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	logger.Fatal(ctx, args...)
}
