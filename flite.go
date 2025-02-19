package flite

import (
	"context"
	"fmt"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
}

type Flite struct {
	Res http.ResponseWriter
	Req *http.Request
	m   *http.ServeMux
	ctx context.Context
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

func (f *Flite) Context(context context.Context) {
	f.ctx = context
}
