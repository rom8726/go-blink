package logs

import (
	"context"
	"os"
)

var EmptyLog Log = emptyLog(0)

type emptyLog int

func (emptyLog) Print(ctx context.Context, l Level, v ...interface{})            {}
func (emptyLog) Printf(ctx context.Context, l Level, f string, v ...interface{}) {}

func (emptyLog) Trace(ctx context.Context, v ...interface{})            {}
func (emptyLog) Tracef(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Debug(ctx context.Context, v ...interface{})            {}
func (emptyLog) Debugf(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Info(ctx context.Context, v ...interface{})             {}
func (emptyLog) Infof(ctx context.Context, f string, v ...interface{})  {}
func (emptyLog) Warn(ctx context.Context, v ...interface{})             {}
func (emptyLog) Warnf(ctx context.Context, f string, v ...interface{})  {}
func (emptyLog) Error(ctx context.Context, v ...interface{})            {}
func (emptyLog) Errorf(ctx context.Context, f string, v ...interface{}) {}
func (emptyLog) Stack(ctx context.Context, v ...interface{})            {}
func (emptyLog) Stackf(ctx context.Context, f string, v ...interface{}) {}

func (emptyLog) Fatal(ctx context.Context, v ...interface{}) {
	os.Exit(1)
}
func (emptyLog) Fatalf(ctx context.Context, f string, v ...interface{}) {
	os.Exit(1)
}
