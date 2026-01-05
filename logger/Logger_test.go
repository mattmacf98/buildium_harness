package logger

import "testing"

// resetSharedLogs clears the global sharedLogs slice between tests
func resetSharedLogs() {
	sharedLogs = nil
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	if logger.step != 0 {
		t.Errorf("NewLogger().step = %d, want 0", logger.step)
	}
}

func TestNextStep(t *testing.T) {
	logger := NewLogger()

	if logger.step != 0 {
		t.Fatalf("initial step = %d, want 0", logger.step)
	}

	logger.NextStep()
	if logger.step != 1 {
		t.Errorf("step after NextStep() = %d, want 1", logger.step)
	}

	logger.NextStep()
	if logger.step != 2 {
		t.Errorf("step after second NextStep() = %d, want 2", logger.step)
	}
}

func TestLogTitle(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 1

	logger.LogTitle("My Test Title")

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Stage != 1 {
		t.Errorf("log.Stage = %d, want 1", log.Stage)
	}
	if log.Message != "My Test Title" {
		t.Errorf("log.Message = %q, want %q", log.Message, "My Test Title")
	}
	if log.Type != "HEADER" {
		t.Errorf("log.Type = %q, want %q", log.Type, "HEADER")
	}
}

func TestLogSuccess(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 2

	logger.LogSuccess("Operation completed")

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Stage != 2 {
		t.Errorf("log.Stage = %d, want 2", log.Stage)
	}
	if log.Message != "Operation completed" {
		t.Errorf("log.Message = %q, want %q", log.Message, "Operation completed")
	}
	if log.Type != "SUCCESS" {
		t.Errorf("log.Type = %q, want %q", log.Type, "SUCCESS")
	}
}

func TestLogInfo(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 3

	logger.LogInfo("Some information")

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Stage != 3 {
		t.Errorf("log.Stage = %d, want 3", log.Stage)
	}
	if log.Message != "Some information" {
		t.Errorf("log.Message = %q, want %q", log.Message, "Some information")
	}
	if log.Type != "INFO" {
		t.Errorf("log.Type = %q, want %q", log.Type, "INFO")
	}
}

func TestLogError(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 4

	logger.LogError("Something went wrong")

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Stage != 4 {
		t.Errorf("log.Stage = %d, want 4", log.Stage)
	}
	if log.Message != "Something went wrong" {
		t.Errorf("log.Message = %q, want %q", log.Message, "Something went wrong")
	}
	if log.Type != "FAILURE" {
		t.Errorf("log.Type = %q, want %q", log.Type, "FAILURE")
	}
}

func TestLogClientCode(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 5

	logger.LogClientCode("fmt.Println(\"hello\")")

	logs := GetAllLogs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Stage != 5 {
		t.Errorf("log.Stage = %d, want 5", log.Stage)
	}
	if log.Message != "fmt.Println(\"hello\")" {
		t.Errorf("log.Message = %q, want %q", log.Message, "fmt.Println(\"hello\")")
	}
	if log.Type != "CLIENT_CODE" {
		t.Errorf("log.Type = %q, want %q", log.Type, "CLIENT_CODE")
	}
}

func TestLogClientCodeMultiline(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 1

	logger.LogClientCode("line1\nline2\nline3")

	logs := GetAllLogs()
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}

	expectedMessages := []string{"line1", "line2", "line3"}
	for i, log := range logs {
		if log.Stage != 1 {
			t.Errorf("logs[%d].Stage = %d, want 1", i, log.Stage)
		}
		if log.Message != expectedMessages[i] {
			t.Errorf("logs[%d].Message = %q, want %q", i, log.Message, expectedMessages[i])
		}
		if log.Type != "CLIENT_CODE" {
			t.Errorf("logs[%d].Type = %q, want %q", i, log.Type, "CLIENT_CODE")
		}
	}
}

func TestLogClientCodeSkipsEmptyLines(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()
	logger.step = 1

	logger.LogClientCode("line1\n\nline2\n\n\nline3")

	logs := GetAllLogs()
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs (empty lines skipped), got %d", len(logs))
	}

	expectedMessages := []string{"line1", "line2", "line3"}
	for i, log := range logs {
		if log.Message != expectedMessages[i] {
			t.Errorf("logs[%d].Message = %q, want %q", i, log.Message, expectedMessages[i])
		}
	}
}

func TestMultipleLogMethods(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()

	logger.NextStep()
	logger.LogTitle("Test 1")
	logger.LogInfo("Starting test")
	logger.LogSuccess("Test passed")

	logger.NextStep()
	logger.LogTitle("Test 2")
	logger.LogError("Test failed")

	logs := GetAllLogs()
	if len(logs) != 5 {
		t.Fatalf("expected 5 logs, got %d", len(logs))
	}

	// Verify stages
	expectedStages := []int{1, 1, 1, 2, 2}
	for i, log := range logs {
		if log.Stage != expectedStages[i] {
			t.Errorf("logs[%d].Stage = %d, want %d", i, log.Stage, expectedStages[i])
		}
	}

	// Verify types
	expectedTypes := []string{"HEADER", "INFO", "SUCCESS", "HEADER", "FAILURE"}
	for i, log := range logs {
		if log.Type != expectedTypes[i] {
			t.Errorf("logs[%d].Type = %q, want %q", i, log.Type, expectedTypes[i])
		}
	}
}

func TestGetAllLogsReturnsCurrentState(t *testing.T) {
	resetSharedLogs()
	logger := NewLogger()

	// Log something
	logger.LogInfo("first")
	logs1 := GetAllLogs()

	if len(logs1) != 1 {
		t.Errorf("logs1 length = %d, want 1", len(logs1))
	}

	// Log more
	logger.LogInfo("second")
	logs2 := GetAllLogs()

	// logs2 should reflect the new state
	if len(logs2) != 2 {
		t.Errorf("logs2 length = %d, want 2", len(logs2))
	}

	// Verify messages
	if logs2[0].Message != "first" {
		t.Errorf("logs2[0].Message = %q, want %q", logs2[0].Message, "first")
	}
	if logs2[1].Message != "second" {
		t.Errorf("logs2[1].Message = %q, want %q", logs2[1].Message, "second")
	}
}
