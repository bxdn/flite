package flite

import (
	"context"
	"net/http"
)

type RequestHandler = func(http.ResponseWriter, *http.Request)

type RequestNode = func(f *flite) (context.Context, error)
