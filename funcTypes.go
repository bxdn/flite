package flite

import (
	"net/http"
)

type RequestHandler = func(http.ResponseWriter, *http.Request)

type RequestNode = func(f *Flite) error
