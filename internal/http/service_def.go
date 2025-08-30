package http

import (
	"github.com/riabininkf/go-modules/di"

	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/log"
)

const DefServiceName = "http.service"

func init() {
	di.Add(
		di.Def[*Service]{
			Name: DefServiceName,
			Build: func(ctn di.Container) (*Service, error) {
				var errorLogger *log.ErrorLogger
				if err := ctn.Fill(log.DefErrorLoggerName, &errorLogger); err != nil {
					return nil, err
				}

				var loginV1 *handlers.LoginV1
				if err := ctn.Fill(handlers.DefLoginV1Name, &loginV1); err != nil {
					return nil, err
				}

				var refreshV1 *handlers.RefreshV1
				if err := ctn.Fill(handlers.DefRefreshV1Name, &refreshV1); err != nil {
					return nil, err
				}

				var registerV1 *handlers.RegisterV1
				if err := ctn.Fill(handlers.DefRegisterV1Name, &registerV1); err != nil {
					return nil, err
				}

				var updatePasswordV1 *handlers.UpdatePasswordV1
				if err := ctn.Fill(handlers.DefUpdatePasswordV1Def, &updatePasswordV1); err != nil {
					return nil, err
				}

				return NewService(
					errorLogger,
					loginV1,
					refreshV1,
					registerV1,
					updatePasswordV1,
				), nil
			},
		},
	)
}
