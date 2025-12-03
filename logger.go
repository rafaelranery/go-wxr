package wxr

import "log"

// Logger defines the interface for logging operations.
// Implementations should handle log messages for debugging and informational purposes.
// A no-op logger is used by default to avoid noisy output in library code.
type Logger interface {
	Printf(format string, v ...any)
}

// noOpLogger is a logger that discards all log messages.
type noOpLogger struct{}

func (n *noOpLogger) Printf(format string, v ...any) {}

// stdLoggerAdapter adapts *log.Logger to the Logger interface.
type stdLoggerAdapter struct {
	logger *log.Logger
}

func (s *stdLoggerAdapter) Printf(format string, v ...any) {
	s.logger.Printf(format, v...)
}
