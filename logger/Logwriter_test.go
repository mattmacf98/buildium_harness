package logger

import (
	"bytes"
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	logger := NewLogger()
	writer := logger.Writer()

	if writer == nil {
		t.Fatal("Writer() returned nil")
	}

	// Verify it implements io.Writer
	var _ io.Writer = writer
}

func TestLogWriter_Write(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 1
	writer := logger.Writer()

	message := []byte("test output from writer")
	n, err := writer.Write(message)

	if err != nil {
		t.Errorf("Write() error = %v, want nil", err)
	}
	if n != len(message) {
		t.Errorf("Write() returned %d, want %d", n, len(message))
	}

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	if logs[0].Type != "CLIENT_CODE" {
		t.Errorf("log.Type = %q, want %q", logs[0].Type, "CLIENT_CODE")
	}
	if logs[0].Message != "test output from writer" {
		t.Errorf("log.Message = %q, want %q", logs[0].Message, "test output from writer")
	}
}

func TestLogWriter_WithIoCopy(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 2
	writer := logger.Writer()

	source := bytes.NewBufferString("copied content")
	_, err := io.Copy(writer, source)

	if err != nil {
		t.Errorf("io.Copy() error = %v, want nil", err)
	}

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	if logs[0].Message != "copied content" {
		t.Errorf("log.Message = %q, want %q", logs[0].Message, "copied content")
	}
}
