package util

import (
	"io"
	"log"
)

// LogLevel to use for a Cluster
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
)

// Logger which allows you to subscribe to specified log events
type Logger struct {
	Level int
	l     *log.Logger
}

// NewLogger creates a new logger with the given module & prefix
func NewLogger(w io.Writer, prefix string) *Logger {
	return &Logger{
		l: log.New(w, prefix, log.LstdFlags),
	}
}

// Debug logs debug output with the debug prefix
func (l *Logger) Debug(m ...interface{}) {
	l.log(LogDebug, "[DEBUG]", m...)
}

// Info logs info output with the info prefix
func (l *Logger) Info(m ...interface{}) {
	l.log(LogInfo, "[INFO]", m...)
}

// Warn logs warn output with the warn prefix
func (l *Logger) Warn(m ...interface{}) {
	l.log(LogWarn, "[WARN]", m...)
}

// Error logs the error output with the error prefix
func (l *Logger) Error(m ...interface{}) {
	l.log(LogError, "[ERROR]", m...)
}

func (l *Logger) log(level int, prefix string, m ...interface{}) {
	if l.Level >= level {
		m = append([]interface{}{prefix}, m...)
		l.l.Println(m...)
	}
}
