package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/riabininkf/go-project-template/internal/logger"
	"github.com/riabininkf/go-project-template/internal/repository/entity"
	"github.com/riabininkf/go-project-template/pb"
)

func (s *server) CreateProductV1(
	ctx context.Context,
	req *pb.CreateProductV1Request,
) (*pb.CreateProductV1Response, error) {
	if req.GetName() == "" {
		s.log.Warn("name is missing")
		return nil, status.New(codes.InvalidArgument, "name is missing").Err()
	}

	var (
		err     error
		isExist bool
	)
	if isExist, err = s.products.ExistsByName(ctx, req.GetName()); err != nil {
		s.log.Error("products.ExistsByName failed", logger.Error(err))
		return nil, ErrInternalServer
	}

	if isExist {
		s.log.Warn("name is busy")
		return nil, status.New(codes.InvalidArgument, "name is busy").Err()
	}

	product := &entity.Product{
		Name:    req.GetName(),
		Comment: sql.NullString{String: req.GetComment(), Valid: req.GetComment() != ""},
	}
	if err = s.products.Add(ctx, product); err != nil {
		s.log.Error("products.Add failed", logger.Error(err))
		return nil, ErrInternalServer
	}

	return &pb.CreateProductV1Response{
		Id:      product.ID,
		Name:    product.Name,
		Comment: product.Comment.String,
	}, nil
}
