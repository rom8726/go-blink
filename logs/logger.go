package logs

import "context"

type LoggerType string

const (
	LoggerTypeDefault LoggerType = ""
	LoggerTypeConsole LoggerType = "console"
	LoggerTypeFile    LoggerType = "file"
)

// logger writes log messages to a destination (a file, a console, etc.)
type logger interface {
	Write(ctx context.Context, message Record)
}

func newLogger(config LoggerConfig) logger {
	switch config.Type {
	case LoggerTypeDefault, LoggerTypeConsole:
		return newConsoleLogger(config)
	case LoggerTypeFile:
		return newFileLogger(config)
	}

	panic("logs: Unsupported logger type " + config.Type)
	return nil
}
