package testcli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
)

type CliTestConfig struct {
	Logger     *logger.Logger
	Executable string
}

func RunCliTest(steps []func(config *CliTestConfig) error, projectId string) {
	path := flag.String("path", "client_bin", "Path to client binary")
	flag.Parse()
	if *path == "" {
		fmt.Println("Path to client binary required")
		os.Exit(1)
	}
	meta := meta.NewMeta(*path)

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps, projectId)
	runner.Run(ctx)
}
