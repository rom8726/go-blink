package logs

import (
	"context"
	"log"
)

// consoleLogger logs messages to the default golang log.
type consoleLogger struct {
	level  Level
	format Format
}

func newConsoleLogger(config LoggerConfig, format Format) *consoleLogger {
	return &consoleLogger{
		level:  config.Level,
		format: format,
	}
}

func (w *consoleLogger) Log(ctx context.Context, msg message) {
	if msg.Level < w.level {
		return
	}

	s := w.format.Format(ctx, msg)
	log.Output(4, s)
}
