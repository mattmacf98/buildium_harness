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

func RunCliTest(steps []func(config *CliTestConfig) error, skipSteps []int) {
	meta := meta.NewMeta()

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps, skipSteps)
	runner.Run(ctx)
}
