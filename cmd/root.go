package main

import (
	"log"
	"os"
	"time"

	"github.com/spec-tacles/spectacles.go/broker"
	"github.com/spec-tacles/spectacles.go/gateway"
	"github.com/spec-tacles/spectacles.go/rest"
	"github.com/spec-tacles/spectacles.go/types"
	"github.com/spf13/cobra"
)

var amqpUrl, amqpGroup, token, logLevel string
var shardCount int

var brokerConnected bool

var logger = log.New(os.Stdout, "[CMD] ", log.Ldate|log.Ltime|log.Lshortfile)
var logLevels = map[string]int{
	"suppress": gateway.LogLevelSuppress,
	"info":     gateway.LogLevelInfo,
	"warn":     gateway.LogLevelWarn,
	"debug":    gateway.LogLevelDebug,
	"error":    gateway.LogLevelError,
}

func main() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "spectacles",
	Short: "Connects to the Discord websocket API using spectacles.go",
	Run: func(cmd *cobra.Command, args []string) {
		amqp := broker.NewAMQP(amqpGroup, "", func(string, []byte) {})
		notSended := make([]*types.ReceivePacket, 0, 100)
		go tryConnect(amqp)

		manager := gateway.NewManager(&gateway.ManagerOptions{
			ShardOptions: &gateway.ShardOptions{
				Identify: &types.Identify{
					Token: token,
				},
			},
			OnPacket: func(shard int, d *types.ReceivePacket) {
				if brokerConnected {
					// send packets that couldn't be sent due to AMQP connection problems.
					if len(notSended) > 0 {
						for _, pk := range notSended {
							amqp.Publish(string(pk.Event), pk.Data)
						}
					}
					amqp.Publish(string(d.Event), d.Data)
				} else {
					notSended = append(notSended, d)
				}
			},
			REST:     rest.NewClient(token),
			LogLevel: logLevels[logLevel],
		})

		if err := manager.Start(); err != nil {
			log.Fatalf("failed to connect to discord: %v", err)
		}
		select {}
	},
}

// tryConnect exponentially increases the retry interval, stopping at 80 seconds
func tryConnect(amqp *broker.AMQP) {
	retryInterval := time.Second * 5
	for err := amqp.Connect(amqpUrl); err != nil; {
		logger.Printf("failed to connect to amqp, retrying in %d seconds: %v\n", retryInterval, err)
		time.Sleep(retryInterval)
		if retryInterval != 80 {
			retryInterval *= 2
		}
	}

	brokerConnected = true
}

func init() {
	rootCmd.Flags().StringVarP(&amqpGroup, "group", "g", "", "The broker group to send Discord events to.")
	rootCmd.Flags().StringVarP(&amqpUrl, "url", "u", "", "The broker URL to connect to.")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "The Discord token used to connect to the gateway.")
	rootCmd.Flags().IntVarP(&shardCount, "shardcount", "c", 0, "The number of shards to spawn.")
	rootCmd.Flags().StringVarP(&logLevel, "loglevel", "l", "info", "The log level.")
}
