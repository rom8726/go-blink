package httpd

import (
	"context"
	"encoding/json"
	"github.com/manucorporat/sse"
	"net/http"
	"sync"
)

type SSEStream struct {
	mu        sync.Mutex
	listeners []SSEStreamListener

	close    chan struct{}
	closed   chan struct{}
	outgoing chan sse.Event

	r *http.Request
	w http.ResponseWriter
}

type SSEStreamListener interface {
	OnSSEStreamOpened(s *SSEStream)
	OnSSEStreamClosed(s *SSEStream)
}

func NewSSEStream(w http.ResponseWriter, r *http.Request) (*SSEStream, error) {
	w.Header().Set("Content-Type", sse.ContentType)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	stream := &SSEStream{
		close:    make(chan struct{}, 1),
		closed:   make(chan struct{}),
		outgoing: make(chan sse.Event, 8),
		r:        r,
		w:        w,
	}
	go stream.loop()
	return stream, nil
}

func (s *SSEStream) Closed() <-chan struct{}    { return s.closed }
func (s *SSEStream) Outgoing() chan<- sse.Event { return s.outgoing }

func (s *SSEStream) Close() <-chan struct{} {
	select {
	case s.close <- struct{}{}:
	default:
	}
	return s.closed
}

func (s *SSEStream) Send(ctx context.Context, text string) {
	event := sse.Event{Data: text}
	select {
	case s.outgoing <- event:
	case <-ctx.Done():
		return
	}
}

func (s *SSEStream) SendJSON(ctx context.Context, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	text := string(bytes)
	s.Send(ctx, text)
	return nil
}

func (s *SSEStream) loop() {
	defer s.onClosed()
	defer close(s.closed)

	for {
		select {
		case <-s.close:
			return
		case <-s.r.Context().Done():
			return
		case event, ok := <-s.outgoing:
			if !ok {
				return
			}

			if err := sse.Encode(s.w, event); err != nil {
				// TODO: Log an error
				return
			}
		}
	}
}

func (s *SSEStream) onClosed() {
	listeners := s.copyListeners()
	for _, l := range listeners {
		l.OnSSEStreamClosed(s)
	}
}

func (s *SSEStream) addListener(l SSEStreamListener) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, listener := range s.listeners {
		if l == listener {
			return
		}
	}
	s.listeners = append(s.listeners, l)

	select {
	case <-s.closed:
		l.OnSSEStreamClosed(s)
	default:
		l.OnSSEStreamOpened(s)
	}
}

func (s *SSEStream) removeListener(l SSEStreamListener) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.listeners[0:0]
	for _, listener := range s.listeners {
		if listener != l {
			filtered = append(filtered, listener)
		}
	}
	s.listeners = filtered
}

func (s *SSEStream) copyListeners() []SSEStreamListener {
	s.mu.Lock()
	if len(s.listeners) == 0 {
		s.mu.Unlock()
		return nil
	}

	listeners := make([]SSEStreamListener, len(s.listeners))
	copy(listeners, s.listeners)
	s.mu.Unlock()
	return listeners
}
