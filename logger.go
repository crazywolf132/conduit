package conduit

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogLevel represents the severity of a log message.
// Valid log levels are LogDebug, LogInfo, LogWarn, and LogError.
type LogLevel int

const (
	// LogDebug is for verbose output, useful during development or troubleshooting.
	LogDebug LogLevel = iota
	// LogInfo is for general informational messages that highlight the progress of the application.
	LogInfo
	// LogWarn is for situations that are unusual or potentially problematic but not yet errors.
	LogWarn
	// LogError is for errors that cause operation failures or other significant problems.
	LogError
)

// Logger is the interface that wraps basic logging methods at various severity levels.
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// DefaultLogger is a simple implementation of Logger that writes to a given io.Writer
// and filters messages based on a minimum LogLevel.
type DefaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger creates a new DefaultLogger with the specified log level and output writer.
// If out is nil, it defaults to os.Stderr.
func NewLogger(level LogLevel, out io.Writer) *DefaultLogger {
	if out == nil {
		out = os.Stderr
	}
	return &DefaultLogger{
		level:  level,
		logger: log.New(out, "", log.LstdFlags),
	}
}

func (l *DefaultLogger) log(level LogLevel, prefix string, v ...interface{}) {
	if level >= l.level {
		l.logger.Print(prefix, " ", fmt.Sprint(v...))
	}
}

func (l *DefaultLogger) logf(level LogLevel, prefix, format string, v ...interface{}) {
	if level >= l.level {
		l.logger.Print(prefix, " ", fmt.Sprintf(format, v...))
	}
}

// Debug logs a message at the Debug level.
func (l *DefaultLogger) Debug(v ...interface{}) { l.log(LogDebug, "[DEBUG]", v...) }

// Info logs a message at the Info level.
func (l *DefaultLogger) Info(v ...interface{}) { l.log(LogInfo, "[INFO]", v...) }

// Warn logs a message at the Warn level.
func (l *DefaultLogger) Warn(v ...interface{}) { l.log(LogWarn, "[WARN]", v...) }

// Error logs a message at the Error level.
func (l *DefaultLogger) Error(v ...interface{}) { l.log(LogError, "[ERROR]", v...) }

// Debugf logs a formatted message at the Debug level.
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.logf(LogDebug, "[DEBUG]", format, v...)
}

// Infof logs a formatted message at the Info level.
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.logf(LogInfo, "[INFO]", format, v...)
}

// Warnf logs a formatted message at the Warn level.
func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	l.logf(LogWarn, "[WARN]", format, v...)
}

// Errorf logs a formatted message at the Error level.
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.logf(LogError, "[ERROR]", format, v...)
}

// NoopLogger is a Logger that discards all log messages.
type NoopLogger struct{}

func (l *NoopLogger) Debug(v ...interface{})                 {}
func (l *NoopLogger) Info(v ...interface{})                  {}
func (l *NoopLogger) Warn(v ...interface{})                  {}
func (l *NoopLogger) Error(v ...interface{})                 {}
func (l *NoopLogger) Debugf(format string, v ...interface{}) {}
func (l *NoopLogger) Infof(format string, v ...interface{})  {}
func (l *NoopLogger) Warnf(format string, v ...interface{})  {}
func (l *NoopLogger) Errorf(format string, v ...interface{}) {}
