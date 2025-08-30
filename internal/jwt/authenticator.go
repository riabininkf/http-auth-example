package jwt

//go:generate mockery --name AccessTokenVerifier --output ./mocks --outpkg mocks --filename access_token_verifier.go --structname AccessTokenVerifier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrTokenMissing = errors.New("jwt token is missing")

func NewAuthenticator(
	verifier AccessTokenVerifier,
	noAuthRoutes []string,
) *Authenticator {
	methods := make(map[string]struct{}, len(noAuthRoutes))
	for _, method := range noAuthRoutes {
		methods[method] = struct{}{}
	}

	return &Authenticator{
		verifier:     verifier,
		noAuthRoutes: methods,
	}
}

type (
	// Authenticator is responsible for JWT authentication and validation of requests.
	Authenticator struct {
		verifier     AccessTokenVerifier
		noAuthRoutes map[string]struct{}
	}

	AccessTokenVerifier interface {
		VerifyAccess(ctx context.Context, token string) (string, error)
	}
)

// Authenticate validates the Authorization header from the request
// and returns the token's subject (user_id) or an error if invalid.
func (a *Authenticator) Authenticate(ctx context.Context, req *http.Request) (string, error) {
	var header string
	if header = req.Header.Get("Authorization"); header == "" || !strings.HasPrefix(header, "Bearer ") {
		if a.isAuthRequired(req) {
			return "", ErrTokenMissing
		}

		return "", nil
	}

	var token string
	if token = strings.TrimSpace(strings.TrimPrefix(header, "Bearer")); token == "" {
		if a.isAuthRequired(req) {
			return "", ErrTokenMissing
		}

		return "", nil
	}

	var (
		err    error
		userID string
	)
	if userID, err = a.verifier.VerifyAccess(ctx, token); err != nil {
		if a.isAuthRequired(req) {
			return "", err
		}

		return "", nil
	}

	return userID, nil
}

// isAuthRequired determines if authentication is required for the given HTTP request
// based on the noAuthRoutes config.
func (a *Authenticator) isAuthRequired(req *http.Request) bool {
	_, ok := a.noAuthRoutes[fmt.Sprintf("%s %s", req.Method, req.URL.Path)]
	return !ok
}
