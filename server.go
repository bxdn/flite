package flite

import (
	"fmt"
	"net/http"
)

type Server interface {
	Register(endpoints ...Endpoint[any])
	Serve(port int) error
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

func (s *server) Register(endpoints ...Endpoint[any]) {
	for _, endpoint := range endpoints {
		s.m.HandleFunc(endpoint.Path(), endpoint.Handler())
	}
}

func (s *server) Serve(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.m)
}