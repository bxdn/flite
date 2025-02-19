package flite

import (
	"fmt"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

type flite struct {
	w http.ResponseWriter
	r *http.Request
	m *http.ServeMux
}

func Flite() *flite {
	f := flite{}
	f.m = http.NewServeMux()
	return &f
}

func (f *flite) CreateEndpoint(path string) *endpoint {
	builder := endpoint{}
	builder.path = path
	builder.f = f
	return &builder
}

func (f *flite) Register(endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		f.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func (f *flite) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), f.m)
}
