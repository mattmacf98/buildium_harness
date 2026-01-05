package testserver

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/mattmacf98/buildium_harness/logger"
	"github.com/mattmacf98/buildium_harness/meta"
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

	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error { return nil },
		func(config *ServerTestConfig) error { return nil },
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

	steps := []func(config *ServerTestConfig) error{}

	runner := NewRunner(m, steps)

	if runner == nil {
		t.Fatal("NewRunner() returned nil")
	}

	if len(runner.steps) != 0 {
		t.Errorf("NewRunner() steps length = %d, want 0", len(runner.steps))
	}
}

func TestNewTestServer(t *testing.T) {
	l := logger.NewLogger()
	server := NewTestServer("/path/to/executable", l)

	if server == nil {
		t.Fatal("NewTestServer() returned nil")
	}

	if server.executable != "/path/to/executable" {
		t.Errorf("NewTestServer() executable = %q, want %q", server.executable, "/path/to/executable")
	}

	if server.logger != l {
		t.Error("NewTestServer() did not set logger correctly")
	}

	if server.running {
		t.Error("NewTestServer() running should be false initially")
	}
}

func TestNewTestServerWithDifferentPaths(t *testing.T) {
	testCases := []struct {
		name       string
		executable string
	}{
		{"simple path", "/app"},
		{"nested path", "/usr/local/bin/myserver"},
		{"relative path", "./bin/server"},
		{"with spaces", "/path/to/my server"},
	}

	l := logger.NewLogger()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(tc.executable, l)
			if server.executable != tc.executable {
				t.Errorf("NewTestServer() executable = %q, want %q", server.executable, tc.executable)
			}
		})
	}
}

func TestServerTestConfigStructFields(t *testing.T) {
	l := logger.NewLogger()
	server := NewTestServer("/path/to/exe", l)

	config := &ServerTestConfig{
		Logger: l,
		Server: server,
	}

	if config.Logger != l {
		t.Error("ServerTestConfig.Logger was not set correctly")
	}

	if config.Server != server {
		t.Error("ServerTestConfig.Server was not set correctly")
	}
}

func TestRunAllStepsPass(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      2,
		Entrypoint: "true", // Use 'true' command which exits immediately
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	callOrder := []int{}
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callOrder = append(callOrder, 0)
			return nil
		},
		func(config *ServerTestConfig) error {
			callOrder = append(callOrder, 1)
			return nil
		},
		func(config *ServerTestConfig) error {
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

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      1, // Only run steps 0 and 1
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
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

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0, // Only run step 0
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
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

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      3,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	expectedError := errors.New("step 1 failed")
	callCount := 0

	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
			callCount++
			return expectedError
		},
		func(config *ServerTestConfig) error {
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

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      2,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	expectedError := errors.New("first step failed")
	callCount := 0

	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callCount++
			return expectedError
		},
		func(config *ServerTestConfig) error {
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

func TestRunConfigHasLogger(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	var receivedLogger *logger.Logger
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
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

func TestRunConfigHasServer(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	var receivedServer *TestServer
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			receivedServer = config.Server
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if receivedServer == nil {
		t.Error("config.Server was nil")
	}
}

func TestRunWithNoSteps(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	m := &meta.Meta{
		Stage:      5,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	steps := []func(config *ServerTestConfig) error{}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() with no steps returned error: %v", err)
	}
}

func TestRunMultipleStepsWithIntermediateFailure(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      4,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	callSequence := []string{}
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callSequence = append(callSequence, "step0")
			return nil
		},
		func(config *ServerTestConfig) error {
			callSequence = append(callSequence, "step1")
			return nil
		},
		func(config *ServerTestConfig) error {
			callSequence = append(callSequence, "step2-fail")
			return errors.New("step 2 failed")
		},
		func(config *ServerTestConfig) error {
			callSequence = append(callSequence, "step3-should-not-run")
			return nil
		},
		func(config *ServerTestConfig) error {
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

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      100, // Much higher than actual steps
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	callCount := 0
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			callCount++
			return nil
		},
		func(config *ServerTestConfig) error {
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

func TestServerStartupTimeEnvVariable(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Test with custom startup time
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "10")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	stepCalled := false
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			stepCalled = true
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if !stepCalled {
		t.Error("Step was not called")
	}
}

func TestServerStartupTimeInvalidEnvVariable(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Test with invalid startup time
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "not-a-number")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err == nil {
		t.Fatal("Run() should have returned an error for invalid SERVER_STARTUP_TIME")
	}
}

func TestServerStartupTimeDefaultValue(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Clear the startup time to test default
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Unsetenv("SERVER_STARTUP_TIME")
	defer func() {
		if originalStartup != "" {
			os.Setenv("SERVER_STARTUP_TIME", originalStartup)
		}
	}()

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	stepCalled := false
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			stepCalled = true
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	// This will use the default 500ms startup time
	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if !stepCalled {
		t.Error("Step was not called with default startup time")
	}
}

func TestTestServerStopWhenNotRunning(t *testing.T) {
	l := logger.NewLogger()
	server := NewTestServer("/usr/bin/true", l)

	// Stop should be safe to call even when not running
	server.Stop()

	if server.running {
		t.Error("Server should not be running after Stop()")
	}
}

func TestTestServerStartAndStop(t *testing.T) {
	l := logger.NewLogger()
	// Use 'sleep' as a simple long-running process
	server := NewTestServer("/bin/sleep", l)

	// Verify server was created and is not running initially
	if server.running {
		t.Error("Server should not be running before Start()")
	}

	// Note: Full Start/Stop testing with actual processes
	// is covered by the integration tests in Run()
}

func TestServerExecutablePath(t *testing.T) {
	// Set ENVIRONMENT to BUILDING to disable supabase calls
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "BUILDING")
	defer os.Setenv("ENVIRONMENT", originalEnv)

	// Use a very short startup time for tests
	originalStartup := os.Getenv("SERVER_STARTUP_TIME")
	os.Setenv("SERVER_STARTUP_TIME", "1")
	defer os.Setenv("SERVER_STARTUP_TIME", originalStartup)

	m := &meta.Meta{
		Stage:      0,
		Entrypoint: "true",
		Path:       "/usr/bin",
		ProjectId:  "test-project-123",
	}

	var serverExecutable string
	steps := []func(config *ServerTestConfig) error{
		func(config *ServerTestConfig) error {
			serverExecutable = config.Server.executable
			return nil
		},
	}

	runner := NewRunner(m, steps)
	ctx := newTestContext()

	err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	expectedPath := "/usr/bin/true"
	if serverExecutable != expectedPath {
		t.Errorf("Server executable = %q, want %q", serverExecutable, expectedPath)
	}
}
