package log

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"
)

const DefErrorLoggerName = "http.log.error_logger"

func init() {
	di.Add(
		di.Def[*ErrorLogger]{
			Name: DefErrorLoggerName,
			Build: func(ctn di.Container) (*ErrorLogger, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				return NewErrorLogger(log), nil
			},
		},
	)
}
