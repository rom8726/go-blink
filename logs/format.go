package logs

import (
	"context"
)

type FormatType string

const (
	FormatDefault FormatType = ""
	FormatString  FormatType = "string"
)

const (
	DefaultFormat     = "${Level}\t${Log}\t${Message}"
	DefaultTimeFormat = "2006/01/02 15:04:05"
)

// format formats a log message into a string.
type format interface {
	Format(ctx context.Context, message message) string
}

func newFormat(config FormatConfig) format {
	switch config.Type {
	case FormatDefault, FormatString:
		return newStringFormat(config)
	}

	panic("logs: Unsupported format type " + config.Type)
	return nil
}

func newDefaultFormat() format {
	return newStringFormat(FormatConfig{
		Message: DefaultFormat,
		Time:    DefaultTimeFormat,
	})
}
