package testserver

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/buildium-org/buildium_harness/logger"
	"github.com/buildium-org/buildium_harness/meta"
	"github.com/buildium-org/buildium_harness/supabase"
)

type Runner struct {
	meta  *meta.Meta
	steps []func(config *ServerTestConfig) error
}

func NewRunner(meta *meta.Meta, steps []func(config *ServerTestConfig) error) *Runner {
	return &Runner{meta: meta, steps: steps}
}

func (r *Runner) Run(ctx context.Context) error {
	l := ctx.Value("logger").(*logger.Logger)
	executable := r.meta.ExecutableDir + "/" + r.meta.Entrypoint
	server := NewTestServer(executable, l)
	ctx = context.WithValue(ctx, "testServer", server)
	supaClient := supabase.NewSupaClient(ctx)
	err := supaClient.Login(ctx)
	if err != nil {
		fmt.Printf("failed to login: %v", err)
		return fmt.Errorf("failed to login: %v", err)
	}
	ctx = context.WithValue(ctx, "supaClient", supaClient)
	completedStage := -1
	for i, step := range r.steps {
		if i > r.meta.Stage {
			break
		}
		err := r.runTest(ctx, step)
		if err != nil {
			supaClient.AddProjectRun(ctx, r.meta.ProjectId, i-1, logger.GetAllLogs())
			return err
		}
		l.NextStep()
		completedStage++
	}
	supaClient.AddProjectRun(ctx, r.meta.ProjectId, completedStage, logger.GetAllLogs())
	return nil
}

func (r *Runner) runTest(ctx context.Context, step func(config *ServerTestConfig) error) error {
	logger := ctx.Value("logger").(*logger.Logger)
	testServer := ctx.Value("testServer").(*TestServer)
	testServer.Start()
	defer testServer.Stop()

	serverStartupTimeStr := os.Getenv("SERVER_STARTUP_TIME")
	var serverStartupTimeMs int = 500
	var err error
	if serverStartupTimeStr != "" {
		serverStartupTimeMs, err = strconv.Atoi(serverStartupTimeStr)
		if err != nil {
			logger.LogError(fmt.Sprintf("invalid server startup time: %v", err))
			return fmt.Errorf("invalid server startup time: %v", err)
		}
	}
	time.Sleep(time.Duration(serverStartupTimeMs) * time.Millisecond)

	err = step(&ServerTestConfig{Logger: logger, Server: testServer})
	if err != nil {
		logger.LogError("Test failed")
		return err
	}
	logger.LogSuccess("Test passed")
	return nil
}
