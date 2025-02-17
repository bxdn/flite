package common

import (
	"context"
	"net/http"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

type RequestNode func(http.ResponseWriter, *http.Request) (context.Context, error)
