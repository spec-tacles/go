package proxy

import (
	"fmt"

	"github.com/spec-tacles/go/broker"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack"
)

type Proxy struct {
	Broker    *broker.AMQP
	RestEvent string
	Token     string
}

func (p *Proxy) Make(method string, path string, body *interface{}, options RequestOptions) (*Response, error) {
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}
	// if they specify Bearer, dont overwrite
	if (options.Headers)["Authorization"] == "" {
		(options.Headers)["Authorization"] = fmt.Sprintf("Bot %s", p.Token)
	}

	if body != nil {
		(options.Headers)["Content-Type"] = "application/json"
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
