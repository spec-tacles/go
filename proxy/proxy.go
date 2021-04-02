package proxy

import (
	"fmt"

	"github.com/spec-tacles/go/broker"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack"
)

type Proxy struct {
	// The AMQP broker used to perform RPC calls to the proxy
	Broker *broker.AMQP
	// The event to call the requests on, usually `REQUEST`
	RestEvent string
	// The Discord bot token to perform requests with
	// To perform Bearer requests, you'll have to specify the header yourself on each DoJSON or Do function
	Token string
}

// generic, private request function
func (p *Proxy) do(method string, path string, body *[]byte, options RequestOptions) (*Response, error) {
	// if they specify Bearer, dont overwrite
	if (options.Headers)["Authorization"] == "" {
		(options.Headers)["Authorization"] = fmt.Sprintf("Bot %s", p.Token)
	}

	data, err := msgpack.Marshal(Request{Method: method, Path: path, Body: body, Headers: options.Headers, Query: options.Query})
	if err != nil {
		return nil, err
	}

	call, err := p.Broker.Call(p.RestEvent, amqp.Publishing{Body: data})
	if err != nil {
		return nil, err
	}

	var res Response
	err = msgpack.Unmarshal(call, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Performs a request to the proxy, specifying the Content-Type header to application/json
func (p *Proxy) DoJSON(method string, path string, body *[]byte, options *RequestOptions) (*Response, error) {
	if options == nil {
		options = newRequestOptions()
	}

	(options.Headers)["Content-Type"] = "application/json"

	return p.do(method, path, body, *options)
}

// Performs a generic request to the proxy, not specifying a Content-Type header
func (p *Proxy) Do(method string, path string, body *[]byte, options *RequestOptions) (*Response, error) {
	if options == nil {
		options = newRequestOptions()
	}

	return p.do(method, path, body, *options)
}
