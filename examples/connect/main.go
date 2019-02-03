package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spec-tacles/spectacles.go/types"

	"github.com/spec-tacles/spectacles.go/gateway"
)

var token = os.Getenv("TOKEN")

type rest struct{}

func (rest) DoJSON(method, path string, body io.Reader, value interface{}) (err error) {
	req, err := http.NewRequest(method, "https://discordapp.com/api"+path, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bot "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(value)
}

func main() {
	c := gateway.NewShard(&gateway.ShardOptions{
		Identify: &types.Identify{
			Token: token,
		},
		OnPacket: func(r *types.ReceivePacket) {
			fmt.Printf("Received op %d, event %s, seq %d\n", r.Op, r.Event, r.Seq)
		},
		LogLevel: gateway.LogLevelDebug,
	})

	var err error
	c.Gateway, err = gateway.FetchGatewayBot(rest{})
	if err != nil {
		log.Panicf("failed to load gateway: %v", err)
	}

	if err := c.Open(); err != nil {
		log.Panicf("failed to open: %v", err)
	}

	select {}
}
