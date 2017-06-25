package logs

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

type Log interface {
	Trace(ctx context.Context, v ...interface{})
	Tracef(ctx context.Context, format string, v ...interface{})

	Debug(ctx context.Context, v ...interface{})
	Debugf(ctx context.Context, format string, v ...interface{})

	Info(ctx context.Context, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})

	Warn(ctx context.Context, v ...interface{})
	Warnf(ctx context.Context, format string, v ...interface{})

	Error(ctx context.Context, v ...interface{})
	Errorf(ctx context.Context, format string, v ...interface{})

	Panic(ctx context.Context, v ...interface{})
	Panicf(ctx context.Context, format string, v ...interface{})

	Fatal(ctx context.Context, v ...interface{})
	Fatalf(ctx context.Context, format string, v ...interface{})
}

func newLog(logs *logs, name string) Log {
	return &logImpl{
		logs: logs,
		name: name,
	}
}

type logImpl struct {
	logs *logs
	name string
}

func (l *logImpl) printf(ctx context.Context, level Level, format string, v ...interface{}) {
	msg := message{
		Log:    l.name,
		Time:   time.Now(),
		Level:  level,
		Format: format,
		Args:   v,
	}
	for _, w := range l.logs.loggers {
		w.Log(ctx, msg)
	}
}

func (l *logImpl) Trace(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelTrace, "", v...)
}

func (l *logImpl) Tracef(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelTrace, f, v...)
}

func (l *logImpl) Debug(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelDebug, "", v...)
}

func (l *logImpl) Debugf(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelDebug, f, v...)
}

func (l *logImpl) Info(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelInfo, "", v...)
}

func (l *logImpl) Infof(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelInfo, f, v...)
}

func (l *logImpl) Warn(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelWarn, "", v...)
}

func (l *logImpl) Warnf(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelWarn, f, v...)
}

func (l *logImpl) Error(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelError, "", v...)
}

func (l *logImpl) Errorf(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelWarn, f, v...)
}

func (l *logImpl) Panic(ctx context.Context, v ...interface{}) {
	msg := fmt.Sprint(v...)
	msg += "\n"
	msg += string(debug.Stack())
	l.printf(ctx, LevelError, "", msg)
}

func (l *logImpl) Panicf(ctx context.Context, f string, v ...interface{}) {
	msg := fmt.Sprintf(f, v...)
	msg += "\n"
	msg += string(debug.Stack())
	l.printf(ctx, LevelError, "", msg)
}

func (l *logImpl) Fatal(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelError, "", v...)
	os.Exit(1)
}

func (l *logImpl) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.printf(ctx, LevelError, format, v...)
	os.Exit(1)
}

type emptyLog int

func (emptyLog) Trace(ctx context.Context, v ...interface{})            {}
func (emptyLog) Tracef(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Debug(ctx context.Context, v ...interface{})            {}
func (emptyLog) Debugf(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Info(ctx context.Context, v ...interface{})             {}
func (emptyLog) Infof(ctx context.Context, f string, v ...interface{})  {}
func (emptyLog) Warn(ctx context.Context, v ...interface{})             {}
func (emptyLog) Warnf(ctx context.Context, f string, v ...interface{})  {}
func (emptyLog) Error(ctx context.Context, v ...interface{})            {}
func (emptyLog) Errorf(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Panic(ctx context.Context, v ...interface{})            {}
func (emptyLog) Panicf(ctx context.Context, f string, v ...interface{}) {}

func (emptyLog) Fatal(ctx context.Context, v ...interface{}) {
	os.Exit(1)
}
func (emptyLog) Fatalf(ctx context.Context, f string, v ...interface{}) {
	os.Exit(1)
}

var EmptyLog Log = emptyLog(0)
