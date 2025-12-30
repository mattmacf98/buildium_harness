package logger

type LogWriter struct {
	logger *Logger
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	l.logger.LogClientCode(string(p))
	return len(p), nil
}
