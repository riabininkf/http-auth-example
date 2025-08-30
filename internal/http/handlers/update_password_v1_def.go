package handlers

import (
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"

	"github.com/riabininkf/http-auth-example/internal/repository"
)

const DefUpdatePasswordV1Def = "http.update-password-v1"

func init() {
	di.Add(
		di.Def[*UpdatePasswordV1]{
			Name: DefUpdatePasswordV1Def,
			Build: func(ctn di.Container) (*UpdatePasswordV1, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				var usersRep *repository.Users
				if err := ctn.Fill(repository.DefUsersName, &usersRep); err != nil {
					return nil, err
				}

				return NewUpdatePasswordV1(
					log,
					usersRep,
					usersRep,
				), nil
			},
		},
	)
}
