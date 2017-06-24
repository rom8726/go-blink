package logs

import (
	"context"
	"fmt"
	"github.com/ivankorobkov/go-blink/str"
	"strings"
)

// stringFormat formats a message into a human readable UTF-8 string using named ${Level} placeholders.
type stringFormat struct {
	f          *str.Formatter
	timeFormat string
	ctx        map[string]string // ${Context} contents represented as map<"Key", "my_key">.
	ctxAttrs   map[string]string // ${Context.Attr} represented as map<"Attr", "Context.Attr">.
}

func newStringFormat(config FormatConfig) Format {
	f := str.NewFormatter(config.Message)

	ctxAttrs := map[string]string{}
	for _, param := range f.Params() {
		if strings.HasPrefix(param, "Context.") {
			key := param[len("Context."):]
			ctxAttrs[key] = param
		}
	}

	timeFormat := config.Time
	if timeFormat == "" {
		timeFormat = DefaultTimeFormat
	}

	return &stringFormat{
		f:          f,
		timeFormat: timeFormat,
		ctx:        config.Context,
		ctxAttrs:   ctxAttrs,
	}
}

func (f *stringFormat) Format(ctx context.Context, msg message) string {
	m := map[string]interface{}{
		"Log":   msg.Log,
		"Time":  msg.Time.Format(f.timeFormat),
		"Level": msg.Level,
	}

	if msg.Format != "" {
		m["Message"] = fmt.Sprintf(msg.Format, msg.Args...)
	} else {
		m["Message"] = fmt.Sprint(msg.Args...)
	}

	if len(f.ctx) > 0 {
		var c map[string]interface{}
		for key, param := range f.ctx {
			val := ctx.Value(key)
			if val == nil {
				continue
			}
			if c == nil {
				c = make(map[string]interface{})
			}
			c[param] = val
		}
		if c != nil {
			m["Context"] = c
		}
	}

	if len(f.ctxAttrs) > 0 {
		for key, param := range f.ctxAttrs {
			val := ctx.Value(key)
			if val == nil {
				continue
			}

			m[param] = val
		}
	}

	return f.f.Format(m)
}
