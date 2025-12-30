package testcli

import (
	"context"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
)

type Runner struct {
	meta  *meta.Meta
	steps []func(config *CliTestConfig) error
}

func NewRunner(meta *meta.Meta, steps []func(config *CliTestConfig) error) *Runner {
	return &Runner{meta: meta, steps: steps}
}

func (r *Runner) Run(ctx context.Context) error {
	logger := ctx.Value("logger").(*logger.Logger)
	executable := r.meta.Path + "/" + r.meta.Entrypoint

	ctx = context.WithValue(ctx, "executable", executable)
	for i, step := range r.steps {
		if i > r.meta.Stage {
			return nil
		}
		err := runTest(ctx, step)
		if err != nil {
			return err
		}
		logger.NextStep()
	}
	return nil
}

func runTest(ctx context.Context, step func(config *CliTestConfig) error) error {
	logger := ctx.Value("logger").(*logger.Logger)
	executable := ctx.Value("executable").(string)

	err := step(&CliTestConfig{Logger: logger, Executable: executable})
	if err != nil {
		logger.LogError("Test failed")
		return err
	}
	logger.Log("Test passed")
	return nil

}
