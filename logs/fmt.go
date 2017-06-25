package logs

import (
	"context"
)

type FormatType string

const (
	FormatDefault FormatType = ""
	FormatString  FormatType = "string"
)

// format formats a Record and returns a record.
type format interface {
	format(ctx context.Context, message Record) string
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
	})
}
