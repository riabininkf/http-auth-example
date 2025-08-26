package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riabininkf/go-modules/db"
	"github.com/riabininkf/go-modules/di"
)

const DefUsersName = "repository.users"

func init() {
	di.Add(
		di.Def[*Users]{
			Name: DefUsersName,
			Build: func(ctn di.Container) (*Users, error) {
				var conn *pgxpool.Pool
				if err := ctn.Fill(db.DefPostgresName, &conn); err != nil {
					return nil, err
				}

				return NewUsers(conn), nil
			},
		},
	)
}
