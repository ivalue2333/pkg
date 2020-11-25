package logx

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	gDefaultLogger = NewDefaultLogger()
)

// NewDefaultLogger
func NewDefaultLogger() Logger {
	formatter := newJSONFormatter()

	l := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	l1 := &DefaultLogger{l}
	SetGLogger(l1)
	return l1
}

func StandardLogger() Logger {
	return gDefaultLogger
}

type DefaultLogger struct {
	l1 *logrus.Logger
}

func (cl *DefaultLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Tracef(format, args...)
}

func (cl *DefaultLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Debugf(format, args...)
}

func (cl *DefaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Infof(format, args...)
}

func (cl *DefaultLogger) Printf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Printf(format, args...)
}

func (cl *DefaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Warnf(format, args...)
}

func (cl *DefaultLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Warningf(format, args...)
}

func (cl *DefaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Errorf(format, args...)
}

func (cl *DefaultLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Fatalf(format, args...)
}

func (cl *DefaultLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	cl.l1.WithContext(ctx).Panicf(format, args...)
}

func (cl *DefaultLogger) Trace(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Trace(args...)
}

func (cl *DefaultLogger) Debug(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Debug(args...)
}

func (cl *DefaultLogger) Info(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Info(args...)
}

func (cl *DefaultLogger) Print(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Print(args...)
}

func (cl *DefaultLogger) Warn(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Warn(args...)
}

func (cl *DefaultLogger) Warning(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Warning(args...)
}

func (cl *DefaultLogger) Error(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Error(args...)
}

func (cl *DefaultLogger) Fatal(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Fatal(args...)
}

func (cl *DefaultLogger) Panic(ctx context.Context, args ...interface{}) {
	cl.l1.WithContext(ctx).Panic(args...)
}
