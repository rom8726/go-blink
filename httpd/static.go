package httpd

import (
	"context"
	"net/http"
)

// NewStaticHandler returns a handler which serves static files from a given directory.
func NewStaticHandler(root http.Dir) Handler {
	fileServer := http.FileServer(root)

	return func(ctx context.Context, req *Req, resp *Resp) error {
		req.URL.Path = req.Param("path")
		fileServer.ServeHTTP(resp, req.Request)
		return nil
	}
}
