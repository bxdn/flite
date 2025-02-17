package flite

type RS interface {
	Endpoints() []Endpoint
}

type Server interface {
	Register(rs RS)
	Serve(port int) error
}
