package logs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStringFormat_Format(t *testing.T) {
	f := newStringFormat(FormatConfig{
		Message: "${Time} ${Level} ${Write} ${Message} ${Context}",
		Context: map[string]string{"ID": "id"},
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "ID", "22fefb70-1cb6-4e3d")

	message := f.format(ctx, Record{
		Log:    "test",
		Time:   time.Time{},
		Level:  LevelInfo,
		Format: "Hello %v",
		Args:   []interface{}{"John Doe"},
	})

	assert.Equal(t, "0001-01-01 00:00:00 INFO test Hello John Doe map[id:22fefb70-1cb6-4e3d]", message)
}
