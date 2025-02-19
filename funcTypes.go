package flite

import (
	"context"
	"net/http"
)

type RequestHandler = func(http.ResponseWriter, *http.Request)

type RequestNode = func(f *Flite) (context.Context, error)
