package handlers

//go:generate mockery --name JwtStorage --output ./mocks --outpkg mocks --filename jwt_storage.go --structname JwtStorage

import "context"

type JwtStorage interface {
	Save(ctx context.Context, token string) error
	Pop(ctx context.Context, token string) error
}
