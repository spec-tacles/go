package gateway

import (
	"log"
	"os"

	"github.com/spec-tacles/spectacles.go/types"
)

// ManagerOptions represents NewManager's options
type ManagerOptions struct {
	ShardOptions *ShardOptions
	REST         REST

	ShardCount  int
	ServerIndex int
	ServerCount int

	OnPacket func(int, *types.ReceivePacket)

	Logger   Logger
	LogLevel int
}

func (opts *ManagerOptions) init() {
	if opts.ServerCount == 0 {
		opts.ServerCount = 1
	}

	if opts.Logger == nil {
		opts.Logger = log.New(os.Stdout, "[Manager] ", log.LstdFlags)
	}
}
