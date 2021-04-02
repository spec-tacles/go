package proxy

// The inner "body" of a proxy response
type ResponseBody struct {
	Status  int32             `msgpack:"status"`
	Headers map[string]string `msgpack:"headers"`
	URL     string            `msgpack:"url"`
	Body    *interface{}      `msgpack:"body"`
}

// The response status
type ResponseStatus int32

const (
	// The request succeeded
	Success ResponseStatus = iota
	// An unknown error occurred
	UnknownError ResponseStatus = iota
	// Invalid request format (non-JSON)
	InvalidRequestFormat ResponseStatus = iota
	// Invalid URL path
	InvalidPath ResponseStatus = iota
	// Invalid URL query
	InvalidQuery ResponseStatus = iota
	// Invalid HTTP method
	InvalidMethod ResponseStatus = iota
	// Invalid headers
	InvalidHeaders ResponseStatus = iota
	// The request failed
	RequestFailure ResponseStatus = iota
	// The request timed out
	RequestTimeout ResponseStatus = iota
)

// A response from the proxy.
// If the Status is 0 (Success), the Body will be a ResponseBody.
// Else, the Body will be of type string, elaborating on the error.
type Response struct {
	Status ResponseStatus `msgpack:"status"`
	Body   interface{}    `msgpack:"body"`
}
