package grpc_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/magiconair/properties/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/riabininkf/go-project-template/internal/grpc"
	"github.com/riabininkf/go-project-template/internal/repository"
	"github.com/riabininkf/go-project-template/internal/repository/entity"
	"github.com/riabininkf/go-project-template/pb"
)

func TestServer_GetProductV1(t *testing.T) {
	testCases := map[string]struct {
		onGet    func() (*entity.Product, error)
		request  *pb.GetProductV1Request
		expError error
	}{
		"product id is missing": {
			request:  &pb.GetProductV1Request{},
			expError: status.New(codes.InvalidArgument, "product id is missing").Err(),
		},
		"products.Get failed": {
			onGet:    func() (*entity.Product, error) { return nil, errors.New("expected error") },
			request:  &pb.GetProductV1Request{Id: gofakeit.Uint64()},
			expError: grpc.ErrInternalServer,
		},
		"product not found": {
			onGet:    func() (*entity.Product, error) { return nil, nil },
			request:  &pb.GetProductV1Request{Id: gofakeit.Uint64()},
			expError: grpc.ErrNotFound,
		},
		"positive case": {
			onGet: func() (*entity.Product, error) {
				return &entity.Product{
					ID:      gofakeit.Uint64(),
					Name:    gofakeit.Name(),
					Comment: sql.NullString{String: gofakeit.Comment(), Valid: true},
				}, nil
			},
			request:  &pb.GetProductV1Request{Id: gofakeit.Uint64()},
			expError: nil,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			products := repository.NewProductsMock(t)

			var product *entity.Product
			if testCase.onGet != nil {
				var err error
				product, err = testCase.onGet()

				products.On("Get", ctx, testCase.request.GetId()).
					Return(product, err)
			}

			server := grpc.NewServer(zap.NewNop(), products)

			resp, err := server.GetProductV1(ctx, testCase.request)
			assert.Equal(t, err, testCase.expError)
			if testCase.expError == nil {
				assert.Equal(t, resp, &pb.GetProductV1Response{
					Id:      product.ID,
					Name:    product.Name,
					Comment: product.Comment.String,
				})
			}
		})
	}
}
