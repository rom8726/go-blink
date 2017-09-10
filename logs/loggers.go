package logs

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

// logger writes log messages to a destination (a file, a console, etc.)
type logger interface {
	log(ctx context.Context, message Record)
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

// consoleLogger logs records to stdout.
type consoleLogger struct {
	level  Level
	format *format
}

func newConsoleLogger(config LoggerConfig) *consoleLogger {
	return &consoleLogger{
		level:  config.Level,
		format: newFormat(config.Message, config.Time, config.Context),
	}
}

func (w *consoleLogger) log(ctx context.Context, record Record) {
	if record.Level < w.level {
		return
	}

	s := w.format.format(ctx, record)
	log.Output(4, s)
}

// fileLogger logs records to files and rotates the files.
type fileLogger struct {
	level  Level
	format *format
	logger *log.Logger
}

func newFileLogger(config LoggerConfig) *fileLogger {
	if err := checkLogFile(config.File); err != nil {
		log.Fatal(err)
	}

	rotated := &lumberjack.Logger{
		Filename:   config.File,
		MaxSize:    config.FileMaxSize,
		MaxAge:     config.FileMaxAge,
		MaxBackups: config.FileMaxBackups,
		LocalTime:  true,
	}

	return &fileLogger{
		level:  config.Level,
		format: newFormat(config.Message, config.Time, config.Context),
		logger: log.New(rotated, "", log.LstdFlags),
	}
}

func (w *fileLogger) log(ctx context.Context, record Record) {
	if record.Level < w.level {
		return
	}

	s := w.format.format(ctx, record)
	w.logger.Output(4, s)
}

func checkLogFile(name string) error {
	// Open or create a file.
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0600)
	if os.IsNotExist(err) {
		file, err = os.Create(name)
	}
	if err != nil {
		return err
	}

	file.Close()
	return nil
}
