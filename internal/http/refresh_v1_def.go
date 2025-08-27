package http

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"

	"github.com/riabininkf/http-auth-example/internal/jwt"
)

const DefRefreshV1Name = "http.refresh-v1"

func init() {
	di.Add(
		di.Def[*httpx.Handler]{
			Name: DefRefreshV1Name,
			Tags: []di.Tag{{Name: TagHandlerName}},
			Build: func(ctn di.Container) (*httpx.Handler, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				var issuer *jwt.Issuer
				if err := ctn.Fill(jwt.DefIssuerName, &issuer); err != nil {
					return nil, err
				}

				var verifier *jwt.Verifier
				if err := ctn.Fill(jwt.DefVerifierName, &verifier); err != nil {
					return nil, err
				}

				return httpx.WrapHandler(log, NewRefreshV1(
					log,
					issuer,
					verifier,
				)), nil
			},
		},
	)
}
