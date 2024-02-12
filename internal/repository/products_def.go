package repository

import (
	"github.com/sarulabs/di/v2"

	"github.com/riabininkf/go-project-template/internal/container"
	"github.com/riabininkf/go-project-template/internal/db"
)

const DefProductsName = "repository.products"

func init() {
	container.Add(di.Def{
		Name: DefProductsName,
		Build: func(ctn di.Container) (interface{}, error) {
			var conn db.Queryer
			if err := container.Fill(db.DefName, &conn); err != nil {
				return nil, err
			}

			return NewProducts(conn), nil
		},
	})
}
