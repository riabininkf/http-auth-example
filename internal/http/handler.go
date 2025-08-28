package http

import "context"

const TagHandlerName = "http.handler"

type TokenIssuer interface {
	IssueAccessToken(userID string) (string, error)
	IssueRefreshToken(userID string) (string, error)
}

type JwtStorage interface {
	Save(ctx context.Context, token string) error
	Pop(ctx context.Context, token string) error
}
