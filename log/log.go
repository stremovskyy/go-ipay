/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
