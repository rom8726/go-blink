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

// Format formats a log message into a string.
type Format interface {
	Format(ctx context.Context, message message) string
}

type FormatConfig struct {
	Name    string
	Type    FormatType
	Message string
	Time    string
	Context map[string]string // Map of context keys to param names
}

func newFormat(config FormatConfig) Format {
	switch config.Type {
	case FormatDefault, FormatString:
		return newStringFormat(config)
	}

	panic("logs: Unsupported format type " + config.Type)
	return nil
}

func newDefaultFormat() Format {
	return newStringFormat(FormatConfig{
		Message: DefaultFormat,
		Time:    DefaultTimeFormat,
	})
}
