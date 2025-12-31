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

type Logger struct {
	step int
	logs []Log
}

func NewLogger() *Logger {
	return &Logger{step: 0}
}

func (l *Logger) Writer() io.Writer {
	return &LogWriter{logger: l}
}

func (l *Logger) NextStep() {
	l.step++
}

func (l *Logger) GetLogs() []Log {
	return l.logs
}

func (l *Logger) LogTitle(title string) {
	fmt.Printf("--------------------------------Test %d: %s--------------------------------\n", l.step, title)
	l.logs = append(l.logs, Log{Stage: l.step, Message: title, Type: "HEADER"})
}

func (l *Logger) Log(message string) {
	fmt.Printf(Colorize(Green, "[Test %d] [Success]: %s\n"), l.step, message)
	l.logs = append(l.logs, Log{Stage: l.step, Message: message, Type: "SUCCESS"})
}

func (l *Logger) LogError(message string) {
	fmt.Printf(Colorize(Red, "[Test %d] [Error]: %s\n"), l.step, message)
	l.logs = append(l.logs, Log{Stage: l.step, Message: message, Type: "FAILURE"})
}

func (l *Logger) LogClientCode(message string) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fmt.Printf(Colorize(Yellow, "[Test %d] [Your Code]: %s\n"), l.step, line)
	}
}
