package log

import "github.com/riabininkf/go-modules/logger"

func NewErrorLogger(
	log *logger.Logger,
) *ErrorLogger {
	return &ErrorLogger{
		log: log,
	}
}

type ErrorLogger struct {
	log *logger.Logger
}

func (l *ErrorLogger) Error(msg string, err error) {
	l.log.Error(msg, logger.Error(err))
}

func (l *ErrorLogger) WithMethod(method string) *ErrorLogger {
	return NewErrorLogger(l.log.With(logger.String("method", method)))
}
