package flite

import (
	"context"
	"fmt"
	"net/http"
)

type Endpoint interface {
	Path() string
	Handler() RequestHandler
	GET(handlers ...RequestNode) Endpoint
	POST(handlers ...RequestNode) Endpoint
	DELETE(handlers ...RequestNode) Endpoint
	PUT(handlers ...RequestNode) Endpoint
}

type Server interface {
	Register(endpoints ...Endpoint)
	Serve(port int) error
}

// Creates an endpoint from a given path.
//
// Uses ServeMux path syntax.
func CreateEndpoint(path string) Endpoint {
	ep := endpoint{}
	ep.path = path
	return &ep
}

type server struct {
	m *http.ServeMux
}

// Creates a server using a new serve mux.
func NewProdServer() Server {
	s := server{http.NewServeMux()}
	return &s
}

// Creates a server using the default serve mux under the hood, exposing the bells and whistles.
// 
// Not recommended for production.
func NewDevServer() Server {
	s := server{http.DefaultServeMux}
	return &s
}

func (s *server) Register(endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		s.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func (s *server) Serve(port int) error {
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
