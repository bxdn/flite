package flite

type Endpoint struct {
	path    string
	handler RequestHandler
}

func (end Endpoint) Path() string {
	return end.path
}

func (end Endpoint) Handler() RequestHandler {
	return end.handler
}
