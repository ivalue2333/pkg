package logx

import "context"

type Logger interface {
	Tracef(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Trace(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Print(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warning(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
}

var gLogger Logger

func SetGLogger(l1 Logger) {
	gLogger = l1
}

func GetGLogger() Logger {
	return gLogger
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
