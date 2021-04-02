// go test proxy/proxy_test.go proxy/proxy.go proxy/response.go proxy/request.go
package proxy

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/spec-tacles/go/broker"
	"github.com/stretchr/testify/assert"
)

var b = broker.NewAMQP("rest", "", nil)
var channel = os.Getenv("CHANNEL_ID")

type PostMessageBody struct {
	Content string `json:"content"`
}

func TestProxy(t *testing.T) {
	err := b.Connect("amqp://admin:appellsmells@localhost//")
	if err != nil {
		panic(err)
	}

	prxy := Proxy{
		Broker:    b,
		RestEvent: "REQUEST",
		Token:     "",
	}

	t.Logf("here")
	res, err := prxy.DoJSON("GET", "/users/@me", nil, nil)
	t.Logf("res: %v", res)
	assert.NoError(t, err)
	assert.EqualValues(t, res.Status, 0)

	body, _ := json.Marshal(PostMessageBody{Content: "foo bar"})

	res, err = prxy.DoJSON("POST", fmt.Sprintf("/channels/%s/messages", channel), &body, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, res.Status, 0)

	switch v := res.Body.(type) {
	case ResponseBody:
		assert.EqualValues(t, v.Status, 200)
	}

	b.Close()
}

// a compose file for your convenience

// version: '3'

// services:
//   rabbit:
//     image: rabbitmq:management
//     restart: unless-stopped
//     environment:
//         - RABBITMQ_DEFAULT_USER=admin
//         - RABBITMQ_DEFAULT_PASS=appellsmells
//     healthcheck:
//         test: ["CMD", "rabbitmq-diagnostics", "-q", "ping"]
//         interval: 10s
//         timeout: 5s
//     ports:
//         - 15672:15672
//         - 5672:5672

//   redis:
//     image: redis:5-alpine
//     restart: unless-stopped
//     volumes:
//       - redis_data:/data
//     healthcheck:
//       test: ['CMD-SHELL', 'redis-cli ping']
//       interval: 10s
//       timeout: 5s
//     expose:
//       - '6379'

//   proxy:
//     image: spectacles/proxy:latest
//     depends_on:
//       rabbit:
//         condition: service_healthy
//     environment:
//       RUST_LOG: trace
//       DISCORD_API_VERSION: 8
//       REDIS_URL: 'redis://redis:6379'
//       AMQP_URL: 'amqp://admin:appellsmells@rabbit//'
//       AMQP_GROUP: 'rest'
//       AMQP_EVENT: 'REQUEST'
//     restart: unless-stopped

// volumes:
//   postgres_data:
//   redis_data:
