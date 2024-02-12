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

func TestServer_CreateProductV1(t *testing.T) {
	testCases := map[string]struct {
		onExistsByName func() (bool, error)
		onAdd          func(product *entity.Product) error
		request        *pb.CreateProductV1Request
		expError       error
	}{
		"name is missing": {
			request:  &pb.CreateProductV1Request{Name: ""},
			expError: status.New(codes.InvalidArgument, "name is missing").Err(),
		},
		"products.ExistsByName failed": {
			onExistsByName: func() (bool, error) { return false, errors.New("expected error") },
			request:        &pb.CreateProductV1Request{Name: gofakeit.Name()},
			expError:       grpc.ErrInternalServer,
		},
		"name is busy": {
			onExistsByName: func() (bool, error) { return true, nil },
			request:        &pb.CreateProductV1Request{Name: gofakeit.Name()},
			expError:       status.New(codes.InvalidArgument, "name is busy").Err(),
		},
		"products.Add failed": {
			onExistsByName: func() (bool, error) { return false, nil },
			onAdd:          func(product *entity.Product) error { return errors.New("expected error") },
			request:        &pb.CreateProductV1Request{Name: gofakeit.Name()},
			expError:       grpc.ErrInternalServer,
		},
		"positive case": {
			onExistsByName: func() (bool, error) { return false, nil },
			onAdd:          func(product *entity.Product) error { return nil },
			request:        &pb.CreateProductV1Request{Name: gofakeit.Name()},
			expError:       nil,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			products := repository.NewProductsMock(t)
			if testCase.onExistsByName != nil {
				products.On("ExistsByName", ctx, testCase.request.GetName()).
					Return(testCase.onExistsByName())
			}

			product := &entity.Product{
				Name: testCase.request.GetName(),
				Comment: sql.NullString{
					String: testCase.request.GetComment(),
					Valid:  testCase.request.GetComment() != "",
				},
			}
			if testCase.onAdd != nil {
				products.On("Add", ctx, product).
					Return(testCase.onAdd(product))
			}

			server := grpc.NewServer(zap.NewNop(), products)

			resp, err := server.CreateProductV1(ctx, testCase.request)
			assert.Equal(t, err, testCase.expError)

			if testCase.expError == nil {
				assert.Equal(t, resp, &pb.CreateProductV1Response{
					Id:      product.ID,
					Name:    product.Name,
					Comment: product.Comment.String,
				})
			}
		})
	}
}
