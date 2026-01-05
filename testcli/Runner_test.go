package testcli

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/buildium-org/buildium_harness/logger"
	"github.com/buildium-org/buildium_harness/meta"
)

// Helper to create a test context with logger
func newTestContext() context.Context {
	l := logger.NewLogger()
	return context.WithValue(context.Background(), "logger", l)
}

func TestNewRunner(t *testing.T) {
	m := &meta.Meta{
		Stage:      2,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error { return nil },
		func(config *CliTestConfig) error { return nil },
	}

	runner := NewRunner(m, steps)

	if runner == nil {
		t.Fatal("NewRunner() returned nil")
	}

	if runner.meta != m {
		t.Error("NewRunner() did not set meta correctly")
	}

	if len(runner.steps) != 2 {
		t.Errorf("NewRunner() steps length = %d, want 2", len(runner.steps))
	}
}

func TestNewRunnerWithEmptySteps(t *testing.T) {
	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	steps := []func(config *CliTestConfig) error{}

	runner := NewRunner(m, steps)

	if runner == nil {
		t.Fatal("NewRunner() returned nil")
	}

	if len(runner.steps) != 0 {
		t.Errorf("NewRunner() steps length = %d, want 0", len(runner.steps))
	}
}

func TestRunAllStepsPass(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      2,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	callOrder := []int{}
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callOrder = append(callOrder, 0)
			return nil
		},
		func(config *CliTestConfig) error {
			callOrder = append(callOrder, 1)
			return nil
		},
		func(config *CliTestConfig) error {
			callOrder = append(callOrder, 2)
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	// All 3 steps should be called (stage 2 means steps 0, 1, 2)
	if len(callOrder) != 3 {
		t.Errorf("Expected 3 steps to be called, got %d", len(callOrder))
	}

	// Verify call order
	for i, v := range callOrder {
		if v != i {
			t.Errorf("Step %d was called in wrong order, got position %d", i, v)
		}
	}
}

func TestRunStopsAtStage(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      1, // Only run steps 0 and 1
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	// Only steps 0 and 1 should be called (stage 1)
	if callCount != 2 {
		t.Errorf("Expected 2 steps to be called for stage 1, got %d", callCount)
	}
}

func TestRunStageZero(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      0, // Only run step 0
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	// Only step 0 should be called
	if callCount != 1 {
		t.Errorf("Expected 1 step to be called for stage 0, got %d", callCount)
	}
}

func TestRunStepFails(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      3,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	expectedError := errors.New("step 1 failed")
	callCount := 0

	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return expectedError
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err == nil {
		t.Fatal("Run() should have returned an error")
	}

	if err != expectedError {
		t.Errorf("Run() error = %v, want %v", err, expectedError)
	}

	// Only steps 0 and 1 should be called (1 fails)
	if callCount != 2 {
		t.Errorf("Expected 2 steps to be called before failure, got %d", callCount)
	}
}

func TestRunFirstStepFails(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      2,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	expectedError := errors.New("first step failed")
	callCount := 0

	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callCount++
			return expectedError
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err == nil {
		t.Fatal("Run() should have returned an error")
	}

	// Only step 0 should be called
	if callCount != 1 {
		t.Errorf("Expected 1 step to be called before failure, got %d", callCount)
	}
}

func TestRunConfigHasCorrectExecutable(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "myapp",
		Path:       "/usr/local/bin",
		ProjectId:  "test-project-123",
	}

	var receivedExecutable string
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			receivedExecutable = config.Executable
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	expectedExecutable := "/usr/local/bin/myapp"
	if receivedExecutable != expectedExecutable {
		t.Errorf("config.Executable = %q, want %q", receivedExecutable, expectedExecutable)
	}
}

func TestRunConfigHasLogger(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	var receivedLogger *logger.Logger
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			receivedLogger = config.Logger
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if receivedLogger == nil {
		t.Error("config.Logger was nil")
	}
}

func TestRunWithNoSteps(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      5,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	steps := []func(config *CliTestConfig) error{}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() with no steps returned error: %v", err)
	}
}

func TestCliTestConfigStructFields(t *testing.T) {
	l := logger.NewLogger()
	config := &CliTestConfig{
		Logger:     l,
		Executable: "/path/to/exe",
	}

	if config.Logger != l {
		t.Error("CliTestConfig.Logger was not set correctly")
	}

	if config.Executable != "/path/to/exe" {
		t.Errorf("CliTestConfig.Executable = %q, want %q", config.Executable, "/path/to/exe")
	}
}

func TestRunMultipleStepsWithIntermediateFailure(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      4,
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	callSequence := []string{}
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callSequence = append(callSequence, "step0")
			return nil
		},
		func(config *CliTestConfig) error {
			callSequence = append(callSequence, "step1")
			return nil
		},
		func(config *CliTestConfig) error {
			callSequence = append(callSequence, "step2-fail")
			return errors.New("step 2 failed")
		},
		func(config *CliTestConfig) error {
			callSequence = append(callSequence, "step3-should-not-run")
			return nil
		},
		func(config *CliTestConfig) error {
			callSequence = append(callSequence, "step4-should-not-run")
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err == nil {
		t.Fatal("Run() should have returned an error")
	}

	expected := []string{"step0", "step1", "step2-fail"}
	if len(callSequence) != len(expected) {
		t.Fatalf("Expected %d calls, got %d: %v", len(expected), len(callSequence), callSequence)
	}

	for i, v := range expected {
		if callSequence[i] != v {
			t.Errorf("callSequence[%d] = %q, want %q", i, callSequence[i], v)
		}
	}
}

func TestRunWithHighStageAndFewerSteps(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      100, // Much higher than actual steps
		Entrypoint: "app",
		Path:       "/test/path",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *CliTestConfig) error{
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
		func(config *CliTestConfig) error {
			callCount++
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	// All available steps should be called
	if callCount != 2 {
		t.Errorf("Expected 2 steps to be called, got %d", callCount)
	}
}
