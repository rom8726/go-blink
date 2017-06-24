package logs

import "context"

type LoggerType string

const (
	LoggerDefault LoggerType = ""
	LoggerConsole LoggerType = "console"
	LoggerFile    LoggerType = "file"
)

// Logger writes log messages to a destination (a file, a console, etc.)
type Logger interface {
	Log(ctx context.Context, message message)
}

type LoggerConfig struct {
	Type   LoggerType
	Level  Level
	Format string

	// File logger
	File           string
	FileMaxSize    int // Maximum size in megabytes of a log file.
	FileMaxAge     int // Maximum number of days to retain old log files.
	FileMaxBackups int // Maximum number of old log files to retain.
}

func newLogger(config LoggerConfig, format Format) Logger {
	switch config.Type {
	case LoggerDefault, LoggerConsole:
		return newConsoleLogger(config, format)
	case LoggerFile:
		return newFileLogger(config, format)
	}

	panic("logs: Unsupported logger type " + config.Type)
	return nil
}
