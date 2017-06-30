package async

import (
	"context"
	"sync"
)

// Service combines a Starter and a Stopper.
type Service interface {
	Starter
	Stopper
}

// Start creates a new service, starts it and returns it.
func Start(run func(ctx context.Context, started chan<- struct{}) error) Service {
	s := NewService(run)
	s.Start()
	return s
}

// NewService creates a new service from a run loop.
// When started, the service creates a goroutine for the run loop.
// The service does not support restarts.
// It is safe to call all the methods in any order.
func NewService(run func(ctx context.Context, started chan<- struct{}) error) Service {
	if run == nil {
		panic("async: nil run")
	}

	return &service{
		run:     run,
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type service struct {
	run func(ctx context.Context, started chan<- struct{}) error
	mu  sync.Mutex

	// Guarded by mu.
	ctx      context.Context    // Context is created when the service is started.
	cancel   context.CancelFunc // Cancel stops the service.
	startErr error
	stopErr  error

	// Goroutine-safe.
	started chan struct{}
	stopped chan struct{}
}

func (s *service) Start() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx == nil {
		s.ctx, s.cancel = context.WithCancel(context.Background())
		go s.main(s.ctx, s.started, s.stopped)
	}

	return s.started
}

func (s *service) Started() <-chan struct{} {
	return s.started
}

func (s *service) StartError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.startErr
}

func (s *service) Stop() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx == nil {
		// Not started.
		s.ctx = context.Background()
		close(s.started)
		close(s.stopped)
	} else {
		s.cancel()
	}

	return s.stopped
}

func (s *service) Stopped() <-chan struct{} {
	return s.stopped
}

func (s *service) StopError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopErr
}

func (s *service) main(ctx context.Context, started chan struct{}, stopped chan<- struct{}) {
	defer close(stopped)
	defer closeOrDefault(started)

	err := s.run(ctx, started)

	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-started:
		s.stopErr = err
	default:
		s.startErr = err
	}
}

func closeOrDefault(ch chan struct{}) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}

// Group groups multiple services into one.
// When a service fails to start, other services which have been started, are stopped.
// The start/stop errors are set to the first errors.
func Group(services ...Service) Service {
	return NewService(func(ctx context.Context, started chan<- struct{}) (err error) {
		defer func() {
			for _, s := range services {
				s.Stop()
			}
			for _, s := range services {
				<-s.Stopped()
			}
			if err != nil {
				return
			}

			for _, s := range services {
				err = s.StopError()
				if err != nil {
					return
				}
			}
		}()

		for _, s := range services {
			s.Start()
		}
		for _, s := range services {
			select {
			case <-s.Started():
			case <-ctx.Done():
				return nil
			}

			err = s.StartError()
			if err != nil {
				return err
			}
		}

		close(started)
		<-ctx.Done()
		return nil
	})
}
