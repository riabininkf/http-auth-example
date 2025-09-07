package jwt

//go:generate mockery --name AccessTokenVerifier --output ./mocks --outpkg mocks --filename access_token_verifier.go --structname AccessTokenVerifier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrTokenMissing indicates that a required JWT token is missing from the request.
var ErrTokenMissing = errors.New("jwt token is missing")

// NewAuthenticator initializes and returns a new instance of Authenticator with the provided verifier and no-auth routes.
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
	// Authenticator handles JWT authentication and manages routes that do not require authentication.
	Authenticator struct {
		verifier     AccessTokenVerifier
		noAuthRoutes map[string]struct{}
	}

	// AccessTokenVerifier defines a method to verify access tokens and return associated user information.
	AccessTokenVerifier interface {
		VerifyAccess(ctx context.Context, token string) (string, error)
	}
)

// Authenticate validates the Authorization header from the HTTP request and extracts the authenticated user ID if valid.
// It uses the provided context and an internal AccessTokenVerifier for token verification.
// Returns the user ID on successful authentication or an error if authentication fails or a token is missing when required.
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

// isAuthRequired checks if authentication is necessary for the given HTTP request based on its method and URL path.
func (a *Authenticator) isAuthRequired(req *http.Request) bool {
	_, ok := a.noAuthRoutes[fmt.Sprintf("%s %s", req.Method, req.URL.Path)]
	return !ok
}
