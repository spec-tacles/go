package proxy

type Request struct {
	_msgpack struct{} `msgpack:",omitempty"`
	// The HTTP method to perform
	Method string `msgpack:"method"`
	// The path to request (e.g: /users/@me)
	Path string `msgpack:"path"`
	// Request headers to include with the request
	Headers map[string]string `msgpack:"headers"`
	// Query parameters to include with the request
	Query map[string]string `msgpack:"query"`
	// The request body
	Body *[]byte `msgpack:"body"`
}

type RequestOptions struct {
	// Headers to include with the request
	Headers map[string]string
	// Query parameters to include with the request
	Query map[string]string
}

func newRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}
}
