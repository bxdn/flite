package flite

import (
	"fmt"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

func Register(endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		http.DefaultServeMux.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), http.DefaultServeMux)
}
