package handlers

//go:generate mockery --name RefreshTokenVerifier --output ./mocks --outpkg mocks --filename refresh_token_verifier.go --structname RefreshTokenVerifier

import (
	"context"
	"net/http"

	"github.com/riabininkf/go-modules/logger"
	"github.com/riabininkf/httpx"
)

// NewRefreshV1 creates a new *RefreshV1 instance.
func NewRefreshV1(
	log *logger.Logger,
	issuer TokenIssuer,
	jwtStorage JwtStorage,
	verifier RefreshTokenVerifier,
) *RefreshV1 {
	return &RefreshV1{
		log:        log,
		issuer:     issuer,
		jwtStorage: jwtStorage,
		verifier:   verifier,
	}
}

type (
	// RefreshV1 generates a new access token using refresh token.
	RefreshV1 struct {
		log        *logger.Logger
		issuer     TokenIssuer
		jwtStorage JwtStorage
		verifier   RefreshTokenVerifier
	}

	// RefreshV1Request represents refresh request.
	RefreshV1Request struct {
		RefreshToken string `json:"refresh_token"`
	}

	// RefreshV1Response represents successful refresh response.
	RefreshV1Response struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// RefreshTokenVerifier describes RefreshTokenVerifier dependency.
	RefreshTokenVerifier interface {
		VerifyRefresh(ctx context.Context, refreshToken string) (string, error)
	}
)

// Handle processes a refresh request, validates the refresh token, generates new tokens, and returns the response.
func (h *RefreshV1) Handle(ctx context.Context, req *RefreshV1Request) *httpx.Response {
	if req.RefreshToken == "" {
		h.log.Warn("refresh_token is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "refresh_token is required")
	}

	if err := h.jwtStorage.Pop(ctx, req.RefreshToken); err != nil {
		h.log.Warn("failed to pop refresh token from the storage", logger.Error(err))
		return httpx.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token")
	}

	var (
		err    error
		userID string
	)
	if userID, err = h.verifier.VerifyRefresh(ctx, req.RefreshToken); err != nil {
		h.log.Warn("failed to verify refresh token", logger.Error(err))
		return httpx.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token")
	}

	var accessToken string
	if accessToken, err = h.issuer.IssueAccessToken(userID); err != nil {
		h.log.Error("failed to issue access token", logger.Error(err))
		return httpx.InternalServerError
	}

	var refreshToken string
	if refreshToken, err = h.issuer.IssueRefreshToken(userID); err != nil {
		h.log.Error("failed to issue refresh token", logger.Error(err))
		return httpx.InternalServerError
	}

	if err = h.jwtStorage.Save(ctx, refreshToken); err != nil {
		h.log.Error("failed to save refresh token", logger.Error(err))
		return httpx.InternalServerError
	}

	return httpx.NewJsonResponse(
		httpx.WithStatus(http.StatusOK),
		httpx.WithBody(&RefreshV1Response{
			UserID:       userID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}),
	)
}
