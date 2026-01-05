package testserver

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/buildium-org/buildium_harness/logger"
	"github.com/buildium-org/buildium_harness/meta"
)

type ServerTestConfig struct {
	Logger *logger.Logger
	Server *TestServer
}

func RunServerTest(steps []func(config *ServerTestConfig) error) {
	path := flag.String("path", "client_bin", "Path to client binary")
	flag.Parse()
	if *path == "" {
		fmt.Println("Path to client binary required")
		os.Exit(1)
	}
	meta := meta.NewMeta(*path)

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps)
	runner.Run(ctx)
}
