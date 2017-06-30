package async

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Start

func TestService_Start__should_start_service(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		close(started)
		<-ctx.Done()
		return nil
	})
	defer s.Stop()
	s.Start()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}
}

func TestService_Start__should_start_service_only_once(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		close(started)
		<-ctx.Done()
		return nil
	})
	defer s.Stop()
	s.Start()
	s.Start()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}
}

func TestService_Start__should_be_ignored_when_already_stopped(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		panic("fatal")
		return nil
	})
	s.Stop()
	s.Start()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}
}

func TestService_StartError__should_return_start_error(t *testing.T) {
	expected := errors.New("test error")
	d := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return expected
	})
	select {
	case <-d.Start():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}

	err := d.StartError()
	assert.Equal(t, expected, err)
}

// Run

func TestService_main__should_mark_service_as_started_on_exit(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return nil
	})
	s.Start()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}
}

func TestService_main__should_set_start_error_on_exit_when_not_started(t *testing.T) {
	testErr := errors.New("test")
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return testErr
	})
	s.Start()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}

	assert.Equal(t, testErr, s.StartError())
	assert.Nil(t, s.StopError())
}

func TestService_main__should_mark_service_as_stopped_on_exit(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return nil
	})
	s.Start()

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}
}

func TestService_main__should_set_stop_error_on_exit_when_started(t *testing.T) {
	testErr := errors.New("test error")
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		close(started)
		return testErr
	})
	s.Start()

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}

	assert.Nil(t, s.StartError())
	assert.Equal(t, testErr, s.StopError())
}

// Stop

func TestService_Stop__should_stop_service__when_started(t *testing.T) {
	d := NewService(func(ctx context.Context, started chan<- struct{}) error {
		close(started)
		<-ctx.Done()
		return nil
	})
	d.Start()
	d.Stop()

	select {
	case <-d.Stopped():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}
}

func TestService_Stop__should_mark_as_started__when_not_started(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return nil
	})
	s.Stop()

	select {
	case <-s.Started():
	case <-time.After(time.Second):
		t.Fatal("Not started")
	}
}

func TestService_Stop__should_mark_as_stopped__when_not_started(t *testing.T) {
	s := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return nil
	})
	s.Stop()

	select {
	case <-s.Stopped():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}
}

func TestService_Stop__should_be_ignored__when_already_exited(t *testing.T) {
	d := NewService(func(ctx context.Context, started chan<- struct{}) error {
		return nil
	})

	select {
	case <-d.Stop():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}

}

func TestService_StopError__should_return_stop_error(t *testing.T) {
	expected := errors.New("Test error")
	d := NewService(func(ctx context.Context, started chan<- struct{}) error {
		close(started)
		<-ctx.Done()
		return expected
	})
	d.Start()

	select {
	case <-d.Stop():
	case <-time.After(time.Second):
		t.Fatal("Not stopped")
	}

	err := d.StopError()
	assert.Equal(t, expected, err)
}
