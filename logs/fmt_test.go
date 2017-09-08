package logs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormat_Format(t *testing.T) {
	f := newFormat(&FormatConfig{
		Message: "${Time} ${Level} ${Log} ${Message} ${Context}",
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

	assert.Equal(t, "0001-01-01 00:00:00 INFO test Hello John Doe id=22fefb70-1cb6-4e3d", message)
}
