package grpc

import (
	"github.com/sarulabs/di/v2"

	"github.com/riabininkf/go-project-template/internal/container"
	"github.com/riabininkf/go-project-template/internal/logger"
	"github.com/riabininkf/go-project-template/internal/repository"
)

const DefName = "grpc.server"

func init() {
	container.Add(di.Def{
		Name: DefName,
		Build: func(ctn di.Container) (interface{}, error) {
			var log logger.Logger
			if err := container.Fill(logger.DefName, &log); err != nil {
				return nil, err
			}

			var products repository.Products
			if err := container.Fill(repository.DefProductsName, &products); err != nil {
				return nil, err
			}

			return NewServer(log, products), nil
		},
	})
}
