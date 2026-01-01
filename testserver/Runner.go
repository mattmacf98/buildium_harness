package testserver

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
	"github.com/mattmacf98/buildium_harness/supabase"
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
	executable := r.meta.Path + "/" + r.meta.Entrypoint
	server := NewTestServer(executable, l)
	ctx = context.WithValue(ctx, "testServer", server)
	email := os.Getenv("BUILDIUM_EMAIL")
	password := os.Getenv("BUILDIUM_PASSWORD")
	if email == "" || password == "" {
		fmt.Printf("BUILDIUM_EMAIL and BUILDIUM_PASSWORD must be set")
		return fmt.Errorf("BUILDIUM_EMAIL and BUILDIUM_PASSWORD must be set")
	}
	supaClient := supabase.NewSupaClient(ctx)
	err := supaClient.Login(ctx, email, password)
	if err != nil {
		fmt.Printf("failed to login: %v", err)
		return fmt.Errorf("failed to login: %v", err)
	}
	ctx = context.WithValue(ctx, "supaClient", supaClient)
	completedStage := 0
	for i, step := range r.steps {
		if i > r.meta.Stage {
			return nil
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
	time.Sleep(500 * time.Millisecond)

	err := step(&ServerTestConfig{Logger: logger, Server: testServer})
	if err != nil {
		logger.LogError("Test failed")
		return err
	}
	logger.Log("Test passed")
	return nil
}
