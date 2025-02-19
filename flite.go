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

func CreateEndpoint(path string) *endpoint {
	ep := endpoint{}
	ep.path = path
	return &ep
}

type Server struct {
	m *http.ServeMux
}

func NewFliteServer() *Server {
	s := Server{}
	s.m = http.NewServeMux()
	return &s
}

func (s *Server) Register(endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		s.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func (s *Server) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.m)
}

type Flite struct {
	Res http.ResponseWriter
	Req *http.Request
	ctx context.Context
}

func (f *Flite) Context(context context.Context) {
	f.ctx = context
}
