package router

import (
	"net/http"
)

type Req struct {
	*http.Request
	Router *Router
	Params Params
}

func newReq(r *Router, r0 *http.Request, params Params) *Req {
	return &Req{
		Router:  r,
		Request: r0,
		Params:  params,
	}
}

func (r *Req) WebSocket(w http.ResponseWriter) (*WebSocket, error) {
	return r.WebSocketOptions(w, nil)
}

func (r *Req) WebSocketOptions(w http.ResponseWriter, opts *WebSocketOptions) (*WebSocket, error) {
	ws, err := NewWebSocket(w, r.Request, opts)
	if err != nil {
		return nil, err
	}

	ws.addListener(r.Router)
	r.Router.onWebSocketOpened(ws)
	return ws, nil
}

func (r *Req) SSE(w http.ResponseWriter) (*SSEStream, error) {
	stream, err := NewSSEStream(w, r.Request)
	if err != nil {
		return nil, err
	}

	stream.addListener(r.Router)
	r.Router.onSSEStreamOpened(stream)
	return stream, nil
}
