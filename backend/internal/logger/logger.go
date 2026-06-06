// Package logger is a tiny zero-dep wrapper around the stdlib log package
// providing levelled output (debug/info/warn/error) and structured key=value
// suffixes. Suitable for production use behind Docker stdout collection.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarn:
		return "WARN "
	case LevelError:
		return "ERROR"
	}
	return "?"
}

func parseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	}
	return LevelInfo
}

// Logger is the package-level logger.
type Logger struct {
	mu    sync.Mutex
	level Level
	out   *log.Logger
}

var std = New(os.Stdout, "info")

// New constructs a logger writing to w with the given minimum level.
func New(w io.Writer, level string) *Logger {
	return &Logger{
		level: parseLevel(level),
		out:   log.New(w, "", log.LstdFlags|log.Lmicroseconds),
	}
}

// SetDefault swaps out the package-level logger.
func SetDefault(l *Logger) { std = l }

// Default returns the package-level logger.
func Default() *Logger { return std }

func (l *Logger) write(lv Level, msg string, kv ...any) {
	if lv < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	suffix := formatKV(kv...)
	if suffix != "" {
		l.out.Printf("%s %s %s", lv, msg, suffix)
	} else {
		l.out.Printf("%s %s", lv, msg)
	}
}

func formatKV(kv ...any) string {
	if len(kv) == 0 {
		return ""
	}
	var parts []string
	for i := 0; i+1 < len(kv); i += 2 {
		parts = append(parts, fmt.Sprintf("%v=%v", kv[i], kv[i+1]))
	}
	return strings.Join(parts, " ")
}

func (l *Logger) Debug(msg string, kv ...any) { l.write(LevelDebug, msg, kv...) }
func (l *Logger) Info(msg string, kv ...any)  { l.write(LevelInfo, msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)  { l.write(LevelWarn, msg, kv...) }
func (l *Logger) Error(msg string, kv ...any) { l.write(LevelError, msg, kv...) }

// Convenience package-level functions.
func Debug(msg string, kv ...any) { std.Debug(msg, kv...) }
func Info(msg string, kv ...any)  { std.Info(msg, kv...) }
func Warn(msg string, kv ...any)  { std.Warn(msg, kv...) }
func Error(msg string, kv ...any) { std.Error(msg, kv...) }
