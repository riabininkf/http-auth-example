package logger

import uberZap "go.uber.org/zap"

func NewZap(logger *uberZap.Logger) Logger {
	return &zap{logger}
}

type (
	zap struct {
		*uberZap.Logger
	}

	Field = uberZap.Field
)

func Error(err error) Field {
	return uberZap.Error(err)
}

func String(key string, val string) Field {
	return uberZap.String(key, val)
}
