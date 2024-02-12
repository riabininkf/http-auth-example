package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/riabininkf/go-project-template/internal/logger"
	"github.com/riabininkf/go-project-template/internal/repository"
	"github.com/riabininkf/go-project-template/pb"
)

var (
	ErrNotFound       = status.New(codes.NotFound, "not found").Err()
	ErrInternalServer = status.New(codes.Internal, "internal server error").Err()
)

func NewServer(
	log logger.Logger,
	products repository.Products,
) pb.TemplateServer {
	return &server{
		log:      log,
		products: products,
	}
}

type (
	Server = pb.TemplateServer

	server struct {
		log      logger.Logger
		products repository.Products
	}
)
