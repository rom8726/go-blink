package logs

import (
	"context"
	"log"
)

// consoleWriter logs records to stdout.
type consoleWriter struct {
	level  Level
	format format
}

func newConsoleWriter(config WriterConfig, format format) *consoleWriter {
	return &consoleWriter{
		level:  config.Level,
		format: format,
	}
}

func (w *consoleWriter) Write(ctx context.Context, record Record) {
	if record.Level < w.level {
		return
	}

	s := w.format.format(ctx, record)
	log.Output(4, s)
}
