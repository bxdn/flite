package server

import (
	"context"
	"net/http"
	"net/url"
)

type WriterFlusher interface {
	http.ResponseWriter
	http.Flusher
}

type statusCacheResponseWriter struct {
	WriterFlusher
	status int
}

func (w *statusCacheResponseWriter) WriteHeader(status int) {
	w.status = status
	w.WriterFlusher.WriteHeader(status)
}

func (w *statusCacheResponseWriter) Write(bytes []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.WriterFlusher.Write(bytes)
}

type F[T any] struct {
	res   *statusCacheResponseWriter
	req   *http.Request
	done  bool
	query url.Values
}

func (f *F[T]) SetContext(context context.Context) {
	f.req = f.req.WithContext(context)
}

func (f *F[T]) AddContext(key, value any) {
	newCtx := context.WithValue(f.req.Context(), key, value)
	f.req = f.req.WithContext(newCtx)
}

func (f *F[T]) Path(key string) string {
	return f.req.PathValue(key)
}

func (f *F[T]) Query(key string) string {
	if f.query == nil {
		f.query = f.req.URL.Query()
	}
	return f.query.Get(key)
}

func (f *F[T]) Header(key string) string {
	return f.req.Header.Get(key)
}

func (f *F[T]) SetHeader(key, value string) {
	f.res.Header().Set(key, value)
}

func (f *F[T]) Req() *http.Request {
	return f.req
}

func (f *F[T]) Res() WriterFlusher {
	return f.res
}
