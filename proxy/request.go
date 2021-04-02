package proxy

type Request struct {
	Method  string             `msgpack:"method"`
	Path    string             `msgpack:"path"`
	Headers map[string]string  `msgpack:"headers"`
	Query   *map[string]string `msgpack:"query"`
	Body    *interface{}       `msgpack:"body"`
}

type RequestOptions struct {
	Headers map[string]string
	Query   *map[string]string
}
