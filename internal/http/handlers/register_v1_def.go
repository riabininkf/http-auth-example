package handlers

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"

	"github.com/riabininkf/http-auth-example/internal/jwt"
	"github.com/riabininkf/http-auth-example/internal/repository"
)

const DefRegisterV1Name = "http.register-v1"

func init() {
	di.Add(
		di.Def[*RegisterV1]{
			Name: DefRegisterV1Name,
			Build: func(ctn di.Container) (*RegisterV1, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				var usersRep *repository.Users
				if err := ctn.Fill(repository.DefUsersName, &usersRep); err != nil {
					return nil, err
				}

				var issuer *jwt.Issuer
				if err := ctn.Fill(jwt.DefIssuerName, &issuer); err != nil {
					return nil, err
				}

				var storage *jwt.Storage
				if err := ctn.Fill(jwt.DefStorageName, &storage); err != nil {
					return nil, err
				}

				return NewRegisterUserV1(
					log,
					issuer,
					storage,
					usersRep,
				), nil
			},
		},
	)
}
