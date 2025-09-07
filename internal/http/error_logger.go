package http

import "github.com/riabininkf/go-modules/logger"

// newErrorLogger creates a new *errorLogger instance.
func newErrorLogger(
	log *logger.Logger,
) *errorLogger {
	return &errorLogger{
		log: log,
	}
}

// errorLogger is an implementation of httpx.Logger with zap.
type errorLogger struct {
	log *logger.Logger
}

// Error implements httpx.Logger.
func (l *errorLogger) Error(msg string, err error) {
	l.log.Error(msg, logger.Error(err))
}
