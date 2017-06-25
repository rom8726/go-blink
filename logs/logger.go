package logs

import "context"

type LoggerType string

const (
	LoggerDefault LoggerType = ""
	LoggerConsole LoggerType = "console"
	LoggerFile    LoggerType = "file"
)

// logger writes log messages to a destination (a file, a console, etc.)
type logger interface {
	Log(ctx context.Context, message message)
}

func newLogger(config LoggerConfig, format format) logger {
	switch config.Type {
	case LoggerDefault, LoggerConsole:
		return newConsoleLogger(config, format)
	case LoggerFile:
		return newFileLogger(config, format)
	}

	panic("logs: Unsupported logger type " + config.Type)
	return nil
}
