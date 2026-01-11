package testserver

import (
	"context"

	"github.com/buildium-org/buildium_harness/logger"
	"github.com/buildium-org/buildium_harness/meta"
)

type ServerTestConfig struct {
	Logger *logger.Logger
	Server *TestServer
}

func RunServerTest(steps []func(config *ServerTestConfig) error, skipSteps []int) {
	meta := meta.NewMeta()

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps, skipSteps)
	runner.Run(ctx)
}
