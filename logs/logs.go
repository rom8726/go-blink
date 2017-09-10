package logs

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	LevelUndefined Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

var levelToName map[Level]string = map[Level]string{
	LevelUndefined: "",
	LevelTrace:     "TRACE",
	LevelDebug:     "DEBUG",
	LevelInfo:      "INFO",
	LevelWarn:      "WARN",
	LevelError:     "ERROR",
}

var nameToLevel map[string]Level = map[string]Level{
	"":      LevelUndefined,
	"TRACE": LevelTrace,
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
}

type LoggerType string

const (
	LoggerTypeDefault LoggerType = ""
	LoggerTypeConsole LoggerType = "console"
	LoggerTypeFile    LoggerType = "file"
)

type Logs interface {
	Log(name string) Log
}

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

type logs struct {
	mu      sync.Mutex
	logs    map[string]Log
	loggers []logger
}

func New(config Config) Logs {
	loggers := []logger{}
	for _, lconf := range config {
		logger := newLogger(lconf)
		loggers = append(loggers, logger)
	}

	return &logs{
		logs:    make(map[string]Log),
		loggers: loggers,
	}
}

func (logs *logs) Log(name string) Log {
	logs.mu.Lock()
	defer logs.mu.Unlock()

	log, ok := logs.logs[name]
	if ok {
		return log
	}

	log = newLog(logs, name)
	logs.logs[name] = log
	return log
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
		Log:     l.name,
		Time:    time.Now(),
		Level:   level,
		Message: format,
		Args:    v,
	}
	for _, logger := range l.logs.loggers {
		logger.log(ctx, record)
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

// Level utility methods

func (level Level) String() string {
	return levelToName[level]
}

func (level Level) MarshalYAML() (interface{}, error) {
	return levelToName[level], nil
}

func (level *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := ""
	if err := unmarshal(&v); err != nil {
		return err
	}
	*level = nameToLevel[strings.ToUpper(v)]
	return nil
}

// UnmarshalJSON implements the json.Marshaler interface.
func (level Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + level.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (level *Level) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.TrimSuffix(s, `"`)
	s = strings.TrimPrefix(s, `"`)
	s = strings.ToUpper(s)
	*level = nameToLevel[s]
	return nil
}
