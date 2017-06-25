package logs

import "context"

type WriterType string

const (
	WriterDefault WriterType = ""
	WriterConsole WriterType = "console"
	WriterFile    WriterType = "file"
)

// writer writes log messages to a destination (a file, a console, etc.)
type writer interface {
	Write(ctx context.Context, message Record)
}

func newWriter(config WriterConfig, format format) writer {
	switch config.Type {
	case WriterDefault, WriterConsole:
		return newConsoleWriter(config, format)
	case WriterFile:
		return newFileWriter(config, format)
	}

	panic("logs: Unsupported writer type " + config.Type)
	return nil
}
