package httpd

import (
	"encoding/json"
	"net/http"
	"net/url"
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

func (r *Req) Int(param string) int      { return r.Params.Int(param) }
func (r *Req) Int32(param string) int32  { return r.Params.Int32(param) }
func (r *Req) Int64(param string) int64  { return r.Params.Int64(param) }
func (r *Req) Param(param string) string { return r.Params[param] }

// Form parses and returns the URL query parameters and the POST data.
// It executes `r.ParseForm(); return r.Request.Form`.
func (r *Req) Form() url.Values {
	r.ParseForm()
	return r.Request.Form
}

// DecodeJSON decodes a JSON body into a given destination.
func (r *Req) DecodeJSON(dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}

func (r *Req) WebSocket(resp *Resp) (*WebSocket, error) {
	return r.WebSocketOptions(resp, nil)
}

func (r *Req) WebSocketOptions(resp *Resp, opts *WebSocketOptions) (*WebSocket, error) {
	ws, err := NewWebSocket(resp.ResponseWriter, r.Request, opts)
	if err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	ws.addListener(r.Router)
	return ws, nil
}

func (r *Req) SSE(resp *Resp) (*SSEStream, error) {
	stream, err := NewSSEStream(resp.ResponseWriter, r.Request)
	if err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	stream.addListener(r.Router)
	return stream, nil
}
