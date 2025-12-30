package testserver

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/mattmacf98/buildium_harness/logger"
)

type TestServer struct {
	executable string
	logger     *logger.Logger
	cleanup    func()
	running    bool
}

func NewTestServer(executable string, logger *logger.Logger) *TestServer {
	return &TestServer{executable: executable, logger: logger}
}

func (t *TestServer) Start() {
	serverCtx := context.Background()
	serverCtx, cancel := context.WithCancel(serverCtx)
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- t.startServer(serverCtx)
	}()

	cleanup := func() {
		cancel()
		<-serverDone // Wait for server to actually terminate
	}
	t.cleanup = cleanup
}

func (t *TestServer) Stop() {
	if !t.running {
		return
	}
	t.cleanup()
}

func (t *TestServer) startServer(ctx context.Context) error {
	cmd := exec.Command(t.executable)
	// Create a new process group so we can kill all child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = t.logger.Writer()
	cmd.Stderr = t.logger.Writer()

	if err := cmd.Start(); err != nil {
		t.logger.LogError(fmt.Sprintf("%v", err))
		return err
	}

	// Wait for context cancellation in a separate goroutine
	go func() {
		<-ctx.Done()
		// Kill the entire process group (negative PID)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		t.running = false
	}()

	t.running = true
	return cmd.Wait()
}
