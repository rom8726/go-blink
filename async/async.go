package async

// Starter is an interface which defines async start and wait methods.
type Starter interface {
	// Start stars a service and returns a channel
	// which is closed when the service starts or fails to start.
	// The method can be called multiple times, but the service starts only once.
	Start() <-chan struct{}

	// Started returns a channel which is closed when the service starts or fails to start.
	Started() <-chan struct{}

	// Returns the start error.
	StartError() error
}

// Stopper is an interface which defines async stop and wait methods.
type Stopper interface {
	// Stop stops a service and returns a channel which is closed when the service stops.
	// The method can be called multiple time.
	Stop() <-chan struct{}

	// Stopped returns a channel which is closed when the service stops.
	Stopped() <-chan struct{}

	// Returns the stop error.
	StopError() error
}
