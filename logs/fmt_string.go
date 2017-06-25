package logs

import (
	"context"
	"fmt"
	"github.com/ivankorobkov/go-blink/strs"
)

// stringFormat formats a Record into a human readable UTF-8 string using named ${Level} placeholders.
type stringFormat struct {
	time string
	fmt  *strs.Formatter
	ctx  map[string]string // ${Context} contents represented as map<"Key", "my_key">.
}

func newStringFormat(config FormatConfig) format {
	timef := config.Time
	if timef == "" {
		timef = DefaultTimeFormat
	}

	return &stringFormat{
		time: timef,
		fmt:  strs.NewFormatter(config.Message),
		ctx:  config.Context,
	}
}

func (f *stringFormat) format(ctx context.Context, record Record) string {
	r := struct {
		Log     string
		Time    string
		Level   Level
		Message string
		Context interface{}
	}{
		Log:     record.Log,
		Time:    record.Time.Format(DefaultTimeFormat),
		Level:   record.Level,
		Message: f.formatMessage(ctx, record),
		Context: f.formatContext(ctx, record),
	}

	return f.fmt.FormatStruct(r)
}

func (f *stringFormat) formatMessage(ctx context.Context, record Record) string {
	if record.Format != "" {
		return fmt.Sprintf(record.Format, record.Args...)
	}
	return fmt.Sprint(record.Args...)
}

func (f *stringFormat) formatContext(ctx context.Context, record Record) interface{} {
	if ctx == nil {
		return ""
	}
	if len(f.ctx) == 0 {
		return ""
	}

	var result map[string]interface{}
	for key, param := range f.ctx {
		val := ctx.Value(key)
		if val == nil {
			continue
		}

		if result == nil {
			result = make(map[string]interface{})
		}
		result[param] = val
	}

	if result == nil {
		return ""
	}
	return result
}
