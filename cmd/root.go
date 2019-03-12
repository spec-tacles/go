package main

import (
	"encoding/json"
	"github.com/spec-tacles/spectacles.go/broker"
	"github.com/spec-tacles/spectacles.go/gateway"
	"github.com/spec-tacles/spectacles.go/rest"
	"github.com/spec-tacles/spectacles.go/types"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var amqpUrl string
var amqpGroup string
var token string
var shardCount int

var logger = log.New(os.Stdout, "[CMD] ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "spectacles",
	Short: "Connects to the Discord websocket API using spectacles.go",
	Run: func(cmd *cobra.Command, args []string) {
		amqp := broker.NewAMQP(amqpGroup, "", func(string, []byte) {})
		done := make(chan bool)
		go tryConnect(amqp, done)

		manager := gateway.NewManager(&gateway.ManagerOptions{
			ShardOptions: &gateway.ShardOptions{
				Identify: &types.Identify{
					Token: token,
				},
			},
			OnPacket: func(shard int, d *types.ReceivePacket) {
				pk, err := json.Marshal(&struct {
					Shard int         `json:"shard_id"`
					Data  interface{} `json:"data"`
				}{shard, d.Data})
				if err != nil {
					return
				}

				select {
				case done <- true:
					amqp.Publish(string(d.Event), pk)
				}
			},
			REST:     rest.NewClient(token),
			LogLevel: gateway.LogLevelDebug,
		})

		if err := manager.Start(); err != nil {
			log.Fatalf("failed to connect to discord: %v", err)
		}
		select {}
	},
}

func tryConnect(amqp *broker.AMQP, done chan bool) {
	for err := amqp.Connect(amqpUrl); err != nil; {
		logger.Printf("failed to connect to amqp, retrying in 30 seconds: %v\n", err)
		time.Sleep(time.Second * 30)
	}

	done <- true
}

func init() {
	rootCmd.Flags().StringVarP(&amqpGroup, "amqpgroup", "g", "", "The AMQP group to send Discord events to.")
	rootCmd.Flags().StringVarP(&amqpUrl, "amqpurl", "u", "", "The AMQP URL to connect to.")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "The Discord token used to connect to the gateway.")
	rootCmd.Flags().IntVarP(&shardCount, "shardcount", "c", 0, "The number of shards to spawn.")
}
