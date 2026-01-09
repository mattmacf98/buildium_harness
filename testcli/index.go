package testcli

import (
	"context"

	"github.com/buildium-org/buildium_harness/logger"
	"github.com/buildium-org/buildium_harness/meta"
)

type CliTestConfig struct {
	Logger     *logger.Logger
	Executable string
}

func RunCliTest(steps []func(config *CliTestConfig) error) {
	meta := meta.NewMetaFromEnv()

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps)
	runner.Run(ctx)
}
