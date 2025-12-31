package testserver

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
)

type ServerTestConfig struct {
	Logger *logger.Logger
	Server *TestServer
}

func RunServerTest(steps []func(config *ServerTestConfig) error) {
	path := flag.String("path", "client_bin", "Path to client binary")
	projectId := flag.String("projectId", "", "Project ID to use for test runs")
	flag.Parse()
	if *path == "" {
		fmt.Println("Path to client binary required")
		os.Exit(1)
	}
	if *projectId == "" {
		fmt.Println("Project ID required")
		os.Exit(1)
	}
	meta := meta.NewMeta(*path)

	logger := logger.NewLogger()
	ctx := context.WithValue(context.Background(), "logger", logger)
	runner := NewRunner(meta, steps, *projectId)
	runner.Run(ctx)
}
