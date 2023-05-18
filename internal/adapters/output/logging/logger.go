// Package logging implements a logger interface
package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// NewSimpleLogger creates a new simpler logger
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

// SimpleLogger implements a simple logger interface
type SimpleLogger struct {
	logger zerolog.Logger
}

// Debug logs a debug message
func (log *SimpleLogger) Debug(msg string, keyvals ...interface{}) error {
	log.logger.Debug().Msgf(msg, keyvals...)
	return nil
}

// Info logs an info message
func (log *SimpleLogger) Info(msg string, keyvals ...interface{}) error {
	log.logger.Info().Msgf(msg, keyvals...)

	return nil
}

// Error logs an error message
func (log *SimpleLogger) Error(msg string, keyvals ...interface{}) error {
	log.logger.Error().Msgf(msg, keyvals...)
	return nil
}
