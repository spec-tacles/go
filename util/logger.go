package util

import (
	"fmt"
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

func NewLogger(w io.Writer, module, prefix string) *Logger {
	return &Logger{
		l: log.New(w, fmt.Sprintf("[%s] %s", module, prefix), log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Debug(m ...interface{}) {
	l.log(LogDebug, "[DEBUG]")
}

func (l *Logger) Info(m ...interface{}) {
	l.log(LogInfo, "[INFO]")
}

func (l *Logger) Warn(m ...interface{}) {
	l.log(LogWarn, "[WARN]")
}

func (l *Logger) Error(m ...interface{}) {
	l.log(LogError, "[ERROR]")
}

func (l *Logger) log(level int, prefix string, m ...interface{}) {
	if l.Level >= level {
		m = append([]interface{}{prefix}, m...)
		l.l.Println(m...)
	}
}
