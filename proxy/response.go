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
	UnknownError
	// Invalid request format (non-JSON)
	InvalidRequestFormat
	// Invalid URL path
	InvalidPath
	// Invalid URL query
	InvalidQuery
	// Invalid HTTP method
	InvalidMethod
	// Invalid headers
	InvalidHeaders
	// The request failed
	RequestFailure
	// The request timed out
	RequestTimeout
)

// A response from the proxy.
// If the Status is 0 (Success), the Body will be a ResponseBody.
// Else, the Body will be of type string, elaborating on the error.
type Response struct {
	Status ResponseStatus `msgpack:"status"`
	Body   interface{}    `msgpack:"body"`
}
