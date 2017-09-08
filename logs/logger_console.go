package logs

import (
	"context"
	"log"
)

// consoleLogger logs records to stdout.
type consoleLogger struct {
	level  Level
	format *format
}

func newConsoleLogger(config LoggerConfig) *consoleLogger {
	return &consoleLogger{
		level:  config.Level,
		format: newFormat(config.Format),
	}
}

func (w *consoleLogger) Write(ctx context.Context, record Record) {
	if record.Level < w.level {
		return
	}

	s := w.format.format(ctx, record)
	log.Output(4, s)
}
