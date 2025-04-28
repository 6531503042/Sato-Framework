package core

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LogLevel represents the log level
type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Fatal
)

// Logger is a configurable logger
type Logger struct {
	level  LogLevel
	output *log.Logger
}

// NewLogger creates a new logger
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		output: log.New(os.Stdout, "", 0),
	}
}

// SetOutput sets the output destination
func (l *Logger) SetOutput(output *log.Logger) {
	l.output = output
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= Debug {
		l.log("DEBUG", format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= Info {
		l.log("INFO", format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= Warn {
		l.log("WARN", format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= Error {
		l.log("ERROR", format, args...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log("FATAL", format, args...)
	os.Exit(1)
}

// log logs a message with the given level
func (l *Logger) log(level string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.output.Printf("[%s] %s %s", timestamp, level, message)
}

// LogMiddleware creates a logging middleware
func LogMiddleware(logger *Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()
		if err != nil {
			logger.Error("Request failed: %v", err)
			return err
		}

		duration := time.Since(start)
		logger.Info("%s %s %d %v", ctx.Method(), ctx.Path(), ctx.Response().StatusCode(), duration)

		return nil
	}
} 