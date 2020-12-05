package logx

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	defaultLogger = NewDefaultLogger()
)

// NewDefaultLogger
func NewDefaultLogger() Logger {
	formatter := newJSONFormatter()

	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	logger := &DefaultLogger{log: log}
	SetLogger(logger)
	return logger
}

func StandardLogger() Logger {
	return defaultLogger
}

type DefaultLogger struct {
	log *logrus.Logger
}

func (cl *DefaultLogger) SetLevel(level logrus.Level) {
	cl.log.SetLevel(level)
}

func (cl *DefaultLogger) SetOutput(out io.Writer) {
	cl.log.SetOutput(out)
}

func (cl *DefaultLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Tracef(format, args...)
}

func (cl *DefaultLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Debugf(format, args...)
}

func (cl *DefaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Infof(format, args...)
}

func (cl *DefaultLogger) Printf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Printf(format, args...)
}

func (cl *DefaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Warnf(format, args...)
}

func (cl *DefaultLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Warningf(format, args...)
}

func (cl *DefaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Errorf(format, args...)
}

func (cl *DefaultLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Fatalf(format, args...)
}

func (cl *DefaultLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	cl.log.WithContext(ctx).Panicf(format, args...)
}

func (cl *DefaultLogger) Trace(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Trace(args...)
}

func (cl *DefaultLogger) Debug(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Debug(args...)
}

func (cl *DefaultLogger) Info(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Info(args...)
}

func (cl *DefaultLogger) Print(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Print(args...)
}

func (cl *DefaultLogger) Warn(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Warn(args...)
}

func (cl *DefaultLogger) Warning(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Warning(args...)
}

func (cl *DefaultLogger) Error(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Error(args...)
}

func (cl *DefaultLogger) Fatal(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Fatal(args...)
}

func (cl *DefaultLogger) Panic(ctx context.Context, args ...interface{}) {
	cl.log.WithContext(ctx).Panic(args...)
}
