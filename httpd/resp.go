package httpd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Resp struct {
	Router *Router
	http.ResponseWriter

	Status     int
	TotalBytes int64
}

func newResp(router *Router, w http.ResponseWriter) *Resp {
	return &Resp{
		Router:         router,
		ResponseWriter: w,
	}
}

func (r *Resp) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.TotalBytes += int64(n)
	return n, err
}

func (r *Resp) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.Status = status
}

func (r *Resp) SetContentType(ctype string) {
	r.Header().Set("Content-Type", ctype)
}

func (r *Resp) SetContentLength(length int64) {
	r.Header().Set("Content-Length", fmt.Sprintf("%d", length))
}

func (r *Resp) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r, cookie)
}

// Text

func (r *Resp) Text(text string) error {
	r.TextStatus(text, http.StatusOK)
	return nil
}

func (r *Resp) TextStatus(text string, status int) error {
	r.SetContentType("text/plain; charset=utf-8")
	r.SetContentLength(int64(len(text)))
	r.WriteHeader(http.StatusOK)
	r.Write([]byte(text))
	return nil
}

// Errors

func (r *Resp) Error(text string, status int) error {
	http.Error(r, text, status)
	return nil
}

func (r *Resp) Error404(text string) error {
	return r.Error(text, http.StatusNotFound)
}

func (r *Resp) Error500(text string) error {
	return r.Error(text, http.StatusInternalServerError)
}

func (r *Resp) ErrorBadRequest() error {
	return r.Error("Bad request", http.StatusBadRequest)
}

func (r *Resp) ErrorInternal() error {
	return r.Error("Internal server error", http.StatusInternalServerError)
}

// JSON

// JSON serves an OK JSON response, returns an error on JSON encoding errors, when no response is written yet.
func (r *Resp) JSON(src interface{}) error {
	return r.JSONStatus(src, http.StatusOK)
}

// JSONStatus returns a JSON response, returns an error on JSON encoding errors, when no response is written yet.
func (r *Resp) JSONStatus(v interface{}, status int) error {
	buf := r.Router.getBuffer()
	defer r.Router.releaseBuffer(buf)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	return r.JSONBytes(buf.Bytes(), status)
}

// JSONBytes serves a JSON response.
func (r *Resp) JSONBytes(bytes []byte, status int) error {
	r.SetContentType("application/json; charset=utf-8")
	r.SetContentLength(int64(len(bytes)))
	r.WriteHeader(http.StatusOK)
	r.Write(bytes)
	return nil
}

// Files

// Attachment serves a file as a downloadable attachment.
func (r *Resp) Attachment(req *Req, name string, modtime time.Time, content io.ReadSeeker) error {
	escaped := url.QueryEscape(name) // The name is already UTF-8.
	disposition := fmt.Sprintf("attachment; filename=%v", escaped)
	r.Header().Set("Content-Disposition", disposition)
	return r.File(req, name, modtime, content)
}

// File serves a file usign http.ServeContent.
func (r *Resp) File(req *Req, name string, modtime time.Time, content io.ReadSeeker) error {
	http.ServeContent(r, req.Request, name, modtime, content)
	return nil
}
