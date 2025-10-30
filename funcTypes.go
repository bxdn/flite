package flite

import (
	"net/http"
)

type RequestHandler = func(http.ResponseWriter, *http.Request)

type RequestNode[T any] func(f *Flite[T]) error
