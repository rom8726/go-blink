package logs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStringFormat_Format(t *testing.T) {
	f := newStringFormat(FormatConfig{
		Message: "${Time} ${Context.ID} ${Level} ${Log} ${Message} ${Context}",
		Context: map[string]string{"UserId": "u"},
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "ID", "22fefb70-1cb6-4e3d")
	ctx = context.WithValue(ctx, "UserId", 1)

	message := f.Format(ctx, message{
		Log:    "test",
		Time:   time.Time{},
		Level:  LevelInfo,
		Format: "Hello %v",
		Args:   []interface{}{"John Doe"},
	})

	assert.Equal(t, "0001/01/01 00:00:00 22fefb70-1cb6-4e3d INFO test Hello John Doe map[u:1]", message)
}
