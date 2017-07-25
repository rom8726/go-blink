package httpd

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

func (r *Req) Int(param string) int {
	return r.Params.Int(param)
}

func (r *Req) Int32(param string) int32 {
	return r.Params.Int32(param)
}

func (r *Req) Int64(param string) int64 {
	return r.Params.Int64(param)
}

func (r *Req) Param(param string) string {
	return r.Params[param]
}

func (r *Req) FormInt(key string) int {
	v := r.FormValue(key)
	i, _ := strconv.ParseInt(v, 10, 64)
	return int(i)
}

func (r *Req) FormInt32(key string) int32 {
	v := r.FormValue(key)
	i, _ := strconv.ParseInt(v, 10, 32)
	return int32(i)
}

func (r *Req) FormInt64(key string) int64 {
	v := r.FormValue(key)
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}

// DecodeJSON decodes a JSON body into a given destination.
func (r *Req) DecodeJSON(dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}

func (r *Req) WebSocket(ctx context.Context, resp *Resp) (*WebSocket, error) {
	ws, err := NewWebSocket(ctx, r.Router.log, resp.ResponseWriter, r.Request)
	if err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	ws.addListener(r.Router)
	return ws, nil
}

func (r *Req) SSEStream(ctx context.Context, resp *Resp) (*SSEStream, error) {
	stream, err := NewSSEStream(ctx, r.Router.log, resp.ResponseWriter, r.Request)
	if err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	stream.addListener(r.Router)
	return stream, nil
}
