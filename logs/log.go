package logs

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

type Log interface {
	Print(ctx context.Context, level Level, v ...interface{})
	Printf(ctx context.Context, level Level, f string, v ...interface{})

	Trace(ctx context.Context, v ...interface{})
	Tracef(ctx context.Context, f string, v ...interface{})

	Debug(ctx context.Context, v ...interface{})
	Debugf(ctx context.Context, f string, v ...interface{})

	Info(ctx context.Context, v ...interface{})
	Infof(ctx context.Context, f string, v ...interface{})

	Warn(ctx context.Context, v ...interface{})
	Warnf(ctx context.Context, f string, v ...interface{})

	Error(ctx context.Context, v ...interface{})
	Errorf(ctx context.Context, f string, v ...interface{})

	Stack(ctx context.Context, v ...interface{})
	Stackf(ctx context.Context, f string, v ...interface{})

	Fatal(ctx context.Context, v ...interface{})
	Fatalf(ctx context.Context, f string, v ...interface{})
}

type logImpl struct {
	logs *logs
	name string
}

func newLog(logs *logs, name string) Log {
	return &logImpl{
		logs: logs,
		name: name,
	}
}

func (l *logImpl) printf(ctx context.Context, level Level, format string, v ...interface{}) {
	record := Record{
		Log:    l.name,
		Time:   time.Now(),
		Level:  level,
		Format: format,
		Args:   v,
	}
	for _, w := range l.logs.loggers {
		w.Write(ctx, record)
	}
}

func (l *logImpl) Print(ctx context.Context, level Level, v ...interface{}) {
	l.printf(ctx, level, "", v...)
}

func (l *logImpl) Printf(ctx context.Context, level Level, format string, v ...interface{}) {
	l.printf(ctx, level, format, v...)
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

func (l *logImpl) Stack(ctx context.Context, v ...interface{}) {
	msg := fmt.Sprint(v...)
	msg += "\n"
	msg += string(debug.Stack())
	l.printf(ctx, LevelError, "", msg)
}

func (l *logImpl) Stackf(ctx context.Context, f string, v ...interface{}) {
	msg := fmt.Sprintf(f, v...)
	msg += "\n"
	msg += string(debug.Stack())
	l.printf(ctx, LevelError, f, msg)
}

func (l *logImpl) Fatal(ctx context.Context, v ...interface{}) {
	l.printf(ctx, LevelError, "", v...)
	os.Exit(1)
}

func (l *logImpl) Fatalf(ctx context.Context, f string, v ...interface{}) {
	l.printf(ctx, LevelError, f, v...)
	os.Exit(1)
}
