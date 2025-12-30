package testserver

import (
	"context"
	"time"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
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
