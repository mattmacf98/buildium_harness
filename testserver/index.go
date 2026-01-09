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

func RunServerTest(steps []func(config *ServerTestConfig) error) {
	meta := meta.NewMetaFromEnv()

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps)
	runner.Run(ctx)
}
