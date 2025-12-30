package logger

import (
	"fmt"
	"io"
	"strings"
)

type Logger struct {
	step int
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

func (l *Logger) LogTitle(title string) {
	fmt.Printf("--------------------------------Test %d: %s--------------------------------\n", l.step, title)
}

func (l *Logger) Log(message string) {
	fmt.Printf(Colorize(Green, "[Test %d] [Success]: %s\n"), l.step, message)
}

func (l *Logger) LogError(message string) {
	fmt.Printf(Colorize(Red, "[Test %d] [Error]: %s\n"), l.step, message)
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
