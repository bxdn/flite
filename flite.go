package flite

import (
	"context"
	"net/http"
)

type statusCacheResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCacheResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusCacheResponseWriter) Write(bytes []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(bytes)
}

type Flite struct {
	res *statusCacheResponseWriter
	req *http.Request
	done bool
}

func (f *Flite) SetContext(context context.Context) {
	f.req = f.req.WithContext(context)
}

func (f *Flite) AddContext(key, value any) {
	newCtx := context.WithValue(f.req.Context(), key, value)
	f.req = f.req.WithContext(newCtx)
}

func (f *Flite) Req() *http.Request{
	return f.req
}

func (f *Flite) Res() http.ResponseWriter {
	return f.res
}
