package logs

import (
	"context"
	"testing"
)

func TestLogs_Log(t *testing.T) {
	logs := New(NewConfig())
	log := logs.Log("test")
	log.Infof(context.Background(), "Hello %v", "world")
}
