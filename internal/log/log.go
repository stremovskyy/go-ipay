package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Level int

const (
	LevelNone    Level = iota // Disables logging.
	LevelError                // Logs anomalies that are not expected to occur during normal use.
	LevelWarning              // Logs anomalies that are expected to occur occasionally during normal use.
	LevelInfo                 // Logs major events.
	LevelDebug                // Logs detailed IO
)

var (
	globalLogLevel Level
	logMutex       sync.Mutex
	labels         = map[Level]string{
		LevelDebug:   "[debug]",
		LevelInfo:    "[info ]",
		LevelWarning: "[warn ]",
		LevelError:   "[error]",
	}
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func SetLevel(level Level) {
	logMutex.Lock()
	defer logMutex.Unlock()
	globalLogLevel = level
}

func (l *Logger) log(level Level, format string, a ...interface{}) {
	if level <= logLevel() {
		prefix := "iPay: "
		if l != nil && l.prefix != "" {
			prefix = l.prefix
		}

		msg := fmt.Sprintf("%s %s %s", time.Now().Format(time.RFC3339), labels[level], prefix)
		msg += fmt.Sprintf(format, a...)
		fmt.Fprintln(os.Stderr, msg)
	}
}

func logLevel() Level {
	logMutex.Lock()
	defer logMutex.Unlock()
	return globalLogLevel
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.log(LevelDebug, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.log(LevelInfo, format, a...)
}

func (l *Logger) Warning(format string, a ...interface{}) {
	l.log(LevelWarning, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.log(LevelError, format, a...)
}
