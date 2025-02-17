package common

type Endpoint struct {
	Path    string
	Handler RequestHandler
}
