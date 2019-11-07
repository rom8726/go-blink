package logs

import (
	"context"
	"fmt"
	"github.com/rom8726/go-blink/strs"
	"strings"
)

const (
	DefaultMessageFormat = "\t${Level}\t${Log}\t${Message}\t${Context}"
	DefaultTimeFormat    = "2006-01-02 15:04:05"
)

// format formats a Record into a human readable UTF-8 string using named ${Level} placeholders.
type format struct {
	formatter   *strs.Formatter
	timeFormat  string
	contextKeys []string
}

func newFormat(messageFormat string, timeFormat string, contextKeys []string) *format {
	if messageFormat == "" {
		messageFormat = DefaultMessageFormat
	}
	if timeFormat == "" {
		timeFormat = DefaultTimeFormat
	}

	return &format{
		formatter:   strs.NewFormatter(messageFormat),
		timeFormat:  timeFormat,
		contextKeys: contextKeys,
	}
}

func (f *format) format(ctx context.Context, record Record) string {
	r := struct {
		Log     string
		Time    string
		Level   Level
		Message string
		Context string
	}{
		Log:     record.Log,
		Time:    record.Time.Format(f.timeFormat),
		Level:   record.Level,
		Message: f.formatMessage(ctx, record),
		Context: f.formatContext(ctx),
	}

	return f.formatter.FormatStruct(r)
}

func (f *format) formatMessage(ctx context.Context, record Record) string {
	if record.Message != "" {
		return fmt.Sprintf(record.Message, record.Args...)
	}
	return fmt.Sprint(record.Args...)
}

func (f *format) formatContext(ctx context.Context) string {
	if ctx == nil || len(f.contextKeys) == 0 {
		return ""
	}

	s := []string{}
	for _, key := range f.contextKeys {
		val := ctx.Value(key)
		if val == nil {
			continue
		}

		s = append(s, fmt.Sprintf("%s=%v", key, val))
	}
	if len(s) == 0 {
		return ""
	}

	return "{" + strings.Join(s, ", ") + "}"
}
