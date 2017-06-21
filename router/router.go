package router

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

var (
	ErrNotFound         = errors.New("router: Not found")
	ErrMethodNotAllowed = errors.New("router: Method not allowed")
)

type Handler func(ctx context.Context, req *Req, resp *Resp) error
type Middleware func(ctx context.Context, req *Req, resp *Resp, next Handler) error

type Router struct {
	NotFoundHandler      http.Handler
	InternalErrorHandler http.Handler

	route      *Route
	streams    map[*SSEStream]struct{}
	websockets map[*WebSocket]struct{}

	mu     sync.Mutex
	close  bool
	closed chan struct{}
}

func New() *Router {
	return &Router{
		route:      NewRoute(),
		streams:    make(map[*SSEStream]struct{}),
		websockets: make(map[*WebSocket]struct{}),
		closed:     make(chan struct{}),
	}
}

func (r *Router) ALL(p string, h Handler)               { r.route.ALL(p, h) }
func (r *Router) HEAD(p string, h Handler)              { r.route.HEAD(p, h) }
func (r *Router) GET(p string, h Handler)               { r.route.GET(p, h) }
func (r *Router) POST(p string, h Handler)              { r.route.POST(p, h) }
func (r *Router) PUT(p string, h Handler)               { r.route.PUT(p, h) }
func (r *Router) DELETE(p string, h Handler)            { r.route.DELETE(p, h) }
func (r *Router) Add(p string, child *Route)            { r.route.Add(p, child) }
func (r *Router) Handler(m string, p string, h Handler) { r.route.Handler(m, p, h) }
func (r *Router) Middleware(p string, m Middleware)     { r.route.Middleware(p, m) }

func (r *Router) ServeHTTP(w http.ResponseWriter, r0 *http.Request) {
	middleware, handler, params, err := r.route.Match(r0.Method, r0.URL.Path)
	if err != nil {
		r.handleError(w, r0, err)
		return
	}

	ctx := r0.Context()
	req := newReq(r, r0, params)
	resp := &Resp{w}
	if err := execute(ctx, middleware, handler, req, resp); err != nil {
		r.handleError(w, r0, err)
	}
}

func (r *Router) handleError(w http.ResponseWriter, r0 *http.Request, err error) {
	switch err {
	case ErrNotFound:
		if r.NotFoundHandler != nil {
			r.NotFoundHandler.ServeHTTP(w, r0)
			return
		}
		http.NotFound(w, r0)

	case ErrMethodNotAllowed:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (r *Router) Close() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.close {
		return r.closed
	}

	r.close = true
	go r.closeAndWait()
	return r.closed
}

func (r *Router) Closed() <-chan struct{} {
	return r.closed
}

func (r *Router) closeAndWait() {
	defer close(r.closed)

	for ws := range r.websockets {
		ws.Close()
	}
	for s := range r.streams {
		s.Close()
	}

	for ws := range r.websockets {
		<-ws.Closed()
	}
	for s := range r.streams {
		<-s.Closed()
	}
}

// sseStreamListener

func (r *Router) onSSEStreamOpened(s *SSEStream) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.close {
		s.Close()
		return
	}

	r.streams[s] = struct{}{}
}

func (r *Router) onSSEStreamClosed(s *SSEStream) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.close {
		return
	}

	delete(r.streams, s)
}

// webSocketListener

func (r *Router) onWebSocketOpened(ws *WebSocket) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.close {
		ws.Close()
		return
	}

	r.websockets[ws] = struct{}{}
}

func (r *Router) onWebSocketClosed(ws *WebSocket) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.close {
		return
	}

	delete(r.websockets, ws)
}

func execute(ctx context.Context, middleware []Middleware, handler Handler, req *Req, resp *Resp) error {
	if len(middleware) == 0 {
		return handler(ctx, req, resp)
	}

	m := middleware[0]
	middleware = middleware[1:]
	return m(ctx, req, resp, func(ctx context.Context, req *Req, resp *Resp) error {
		return execute(ctx, middleware, handler, req, resp)
	})
}
