package main

import (
	"encoding/json"

	"github.com/spec-tacles/go/broker"
	"github.com/streadway/amqp"
)

type ProxyData struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Query   *string           `json:"query"`
	Body    *interface{}      `json:"body"`
	Headers map[string]string `json:"headers"`
}

var b = broker.NewAMQP("rest", "", nil)

func main() {

	err := b.Connect("amqp://admin:doctordoctor@localhost//")
	if err != nil {
		panic(err)
	}

	data := ProxyData{
		Method: "GET",
		Path:   "/users/@me",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bot NzA3OTk1NjY1NjYxMTY1NTY4.XrQ6WA.gxPDv5knXPUFWTLOp4yavB58jFE",
		},
	}
	body, err := json.Marshal(data)
	println(string(body))
	if err != nil {
		panic(err)
	}

	res, err := b.Call("REQUEST", amqp.Publishing{Body: body})
	if err != nil {
		panic(err)
	}
	println(string(res))
	var x interface{}
	err = json.Unmarshal(res, &x)
	println(x)
}
