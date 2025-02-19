package flite

import (
	"fmt"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

type Flite struct {
	w http.ResponseWriter
	r *http.Request
	m *http.ServeMux
}

func NewFlite() *Flite {
	f := Flite{}
	f.m = http.NewServeMux()
	return &f
}

func (f *Flite) CreateEndpoint(path string) *endpoint {
	builder := endpoint{}
	builder.path = path
	builder.f = f
	return &builder
}

func (f *Flite) Register(endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		f.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func (f *Flite) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), f.m)
}
