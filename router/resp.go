package router

import "net/http"

type Resp struct {
	http.ResponseWriter
}

func newResp(w http.ResponseWriter) *Resp {
	return &Resp{w}
}
