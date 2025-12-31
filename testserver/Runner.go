package testserver

import (
	"context"
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
	logger := ctx.Value("logger").(*logger.Logger)
	executable := r.meta.Path + "/" + r.meta.Entrypoint
	server := NewTestServer(executable, logger)
	ctx = context.WithValue(ctx, "testServer", server)
	supaClient := supabase.NewSupaClient(ctx)
	ctx = context.WithValue(ctx, "supaClient", supaClient)
	for i, step := range r.steps {
		if i > r.meta.Stage {
			return nil
		}
		err := r.runTest(ctx, step)
		if err != nil {
			return err
		}
		logger.NextStep()
	}
	return nil
}

func (r *Runner) runTest(ctx context.Context, step func(config *ServerTestConfig) error) error {
	logger := ctx.Value("logger").(*logger.Logger)
	testServer := ctx.Value("testServer").(*TestServer)
	supaClient := ctx.Value("supaClient").(*supabase.SupaClient)
	testServer.Start()
	defer testServer.Stop()
	time.Sleep(500 * time.Millisecond)

	err := step(&ServerTestConfig{Logger: logger, Server: testServer})
	if err != nil {
		logger.LogError("Test failed")
		supaClient.CallDemoFunction(ctx, "FAILED")
		return err
	}
	logger.Log("Test passed")
	supaClient.CallDemoFunction(ctx, "PASSED")
	return nil
}
