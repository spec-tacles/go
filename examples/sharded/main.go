package main

import (
	"encoding/json"
	"errors"
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
	m := gateway.NewManager(&gateway.ManagerOptions{
		ShardOptions: &gateway.ShardOptions{
			Identify: &types.Identify{
				Token: token,
			},
			LogLevel: gateway.LogLevelDebug,
		},
		REST: rest{},
		OnPacket: func(shard int, r *types.ReceivePacket) {
			// if r.Event == "PRESENCE_UPDATE" || r.Event == "MESSAGE_CREATE" || r.Event == "TYPING_START" || r.Event == "GUILD_CREATE" {
			// 	return
			// }

			// fmt.Printf("Received op %d, event %s, seq %d on shard %d\n", r.Op, r.Event, r.Seq, shard)
		},
		LogLevel: gateway.LogLevelInfo,
	})

	if err := m.Start(); err != nil {
		log.Panicf("failed to start: %v", err)
	}

	select {}
}
