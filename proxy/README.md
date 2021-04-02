# go/proxy
a client for using https://github.com/spec-tacles/proxy

## example
```go
import (
	"github.com/spec-tacles/go/broker"
	"github.com/spec-tacles/go/proxy"
)

// ...
b := broker.NewAMQP("rest", "", nil)
// conect broker

prxy := Proxy{
	Broker:    b,
	RestEvent: "REQUEST",
	Token:     "",
}

res, err := prxy.Make("GET", "/users/@me", nil, RequestOptions{})
```