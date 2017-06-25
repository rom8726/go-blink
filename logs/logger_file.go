package logs

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

// fileLogger logs messages to files and rotates the files.
type fileLogger struct {
	level  Level
	format format
	logger *log.Logger
}

func newFileLogger(config LoggerConfig, format format) *fileLogger {
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
		format: format,
		logger: log.New(rotated, "", log.LstdFlags),
	}
}

func (w *fileLogger) Log(ctx context.Context, msg message) {
	if msg.Level < w.level {
		return
	}

	s := w.format.Format(ctx, msg)
	log.Output(4, s)
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
