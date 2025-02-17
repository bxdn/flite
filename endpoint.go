package flite

type Endpoint struct {
	Path    string
	Handler RequestHandler
}
