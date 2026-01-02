package logger

import (
	"fmt"
	"io"
	"strings"
)

type Log struct {
	Stage   int    `json:"stage"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

var sharedLogs []Log

type Logger struct {
	step int
}

func NewLogger() *Logger {
	return &Logger{step: 0}
}

func GetAllLogs() []Log {
	return sharedLogs
}

func (l *Logger) Writer() io.Writer {
	return &LogWriter{logger: l}
}

func (l *Logger) NextStep() {
	l.step++
}

func (l *Logger) LogTitle(title string) {
	fmt.Printf("--------------------------------Test %d: %s--------------------------------\n", l.step, title)
	sharedLogs = append(sharedLogs, Log{Stage: l.step, Message: title, Type: "HEADER"})
}

func (l *Logger) LogSuccess(message string) {
	fmt.Printf(Colorize(Green, "[Test %d] [Success]: %s\n"), l.step, message)
	sharedLogs = append(sharedLogs, Log{Stage: l.step, Message: message, Type: "SUCCESS"})
}

func (l *Logger) LogInfo(message string) {
	fmt.Printf(Colorize(Blue, "[Test %d] [Info]: %s\n"), l.step, message)
	sharedLogs = append(sharedLogs, Log{Stage: l.step, Message: message, Type: "INFO"})
}

func (l *Logger) LogError(message string) {
	fmt.Printf(Colorize(Red, "[Test %d] [Error]: %s\n"), l.step, message)
	sharedLogs = append(sharedLogs, Log{Stage: l.step, Message: message, Type: "FAILURE"})
}

func (l *Logger) LogClientCode(message string) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fmt.Printf(Colorize(Yellow, "[Test %d] [Your Code]: %s\n"), l.step, line)
		sharedLogs = append(sharedLogs, Log{Stage: l.step, Message: line, Type: "CLIENT_CODE"})
	}
}
