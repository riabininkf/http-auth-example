package http

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"

	"github.com/riabininkf/http-auth-example/internal/auth"
	"github.com/riabininkf/http-auth-example/internal/repository"
)

const DefLoginV1Name = "http.login-v1"

func init() {
	di.Add(
		di.Def[*httpx.Handler]{
			Name: DefLoginV1Name,
			Tags: []di.Tag{{Name: TagHandlerName}},
			Build: func(ctn di.Container) (*httpx.Handler, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				var issuer *auth.Issuer
				if err := ctn.Fill(auth.DefIssuerName, &issuer); err != nil {
					return nil, err
				}

				var usersRep *repository.Users
				if err := ctn.Fill(repository.DefUsersName, &usersRep); err != nil {
					return nil, err
				}

				return httpx.WrapHandler(log, NewLogiV1(
					log,
					issuer,
					usersRep,
				)), nil
			},
		},
	)
}
