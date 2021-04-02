package proxy

type Request struct {
	_msgpack struct{}          `msgpack:",omitempty"`
	Method   string            `msgpack:"method"`
	Path     string            `msgpack:"path"`
	Headers  map[string]string `msgpack:"headers"`
	Query    map[string]string `msgpack:"query"`
	Body     *[]byte           `msgpack:"body"`
}

type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
}

func newRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}
}
