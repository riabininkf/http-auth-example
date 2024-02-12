package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/riabininkf/go-project-template/internal/logger"
	"github.com/riabininkf/go-project-template/internal/repository/entity"
	"github.com/riabininkf/go-project-template/pb"
)

func (s *server) GetProductV1(
	ctx context.Context,
	req *pb.GetProductV1Request,
) (*pb.GetProductV1Response, error) {
	if req.GetId() < 1 {
		s.log.Warn("product id is missing")
		return nil, status.New(codes.InvalidArgument, "product id is missing").Err()
	}

	var (
		err     error
		product *entity.Product
	)
	if product, err = s.products.Get(ctx, req.GetId()); err != nil {
		s.log.Error("products.Get failed", logger.Error(err))
		return nil, ErrInternalServer
	}

	if product == nil {
		s.log.Warn("product not found")
		return nil, ErrNotFound
	}

	return &pb.GetProductV1Response{
		Id:      product.ID,
		Name:    product.Name,
		Comment: product.Comment.String,
	}, nil
}
