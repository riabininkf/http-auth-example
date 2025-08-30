package handlers

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"

	"github.com/riabininkf/http-auth-example/internal/jwt"
)

const DefRefreshV1Name = "http.refresh-v1"

func init() {
	di.Add(
		di.Def[*RefreshV1]{
			Name: DefRefreshV1Name,
			Build: func(ctn di.Container) (*RefreshV1, error) {
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

				var storage *jwt.Storage
				if err := ctn.Fill(jwt.DefStorageName, &storage); err != nil {
					return nil, err
				}

				return NewRefreshV1(
					log,
					issuer,
					storage,
					verifier,
				), nil
			},
		},
	)
}
