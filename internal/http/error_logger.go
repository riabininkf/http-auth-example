package http

import "github.com/riabininkf/go-modules/logger"

func newErrorLogger(
	log *logger.Logger,
) *errorLogger {
	return &errorLogger{
		log: log,
	}
}

type errorLogger struct {
	log *logger.Logger
}

func (l *errorLogger) Error(msg string, err error) {
	l.log.Error(msg, logger.Error(err))
}
