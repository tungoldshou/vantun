package core

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger with configurable log level
type Logger struct {
	level LogLevel
	mu    sync.Mutex
}

// NewLogger creates a new logger with the specified log level
func NewLogger(level string) *Logger {
	l := &Logger{}
	
	switch level {
	case "debug":
		l.level = DebugLevel
	case "info":
		l.level = InfoLevel
	case "warn":
		l.level = WarnLevel
	case "error":
		l.level = ErrorLevel
	default:
		l.level = InfoLevel
	}
	
	return l
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	switch level {
	case "debug":
		l.level = DebugLevel
	case "info":
		l.level = InfoLevel
	case "warn":
		l.level = WarnLevel
	case "error":
		l.level = ErrorLevel
	default:
		l.level = InfoLevel
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, format, args...)
	}
}

// log writes a log message with the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, args...)
	
	// Write to stderr
	log.Printf("[%s] %s: %s\n", timestamp, level.String(), message)
	
	// Also write to a log file if needed
	// This is a simplified implementation - in a production system,
	// you might want to use a more robust logging library like logrus or zap
}

// Global logger instance
var defaultLogger *Logger

// InitLogger initializes the global logger with the specified log level
func InitLogger(level string) {
	defaultLogger = NewLogger(level)
}

// Debug logs a debug message using the global logger
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

// Info logs an info message using the global logger
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

// Warn logs a warning message using the global logger
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

// Error logs an error message using the global logger
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}