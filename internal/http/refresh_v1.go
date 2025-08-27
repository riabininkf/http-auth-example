package http

import (
	"context"
	"net/http"

	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"
)

func NewRefreshV1(
	log *logger.Logger,
	issuer TokenIssuer,
	verifier RefreshTokenVerifier,
) *RefreshV1 {
	return &RefreshV1{
		log:      log,
		issuer:   issuer,
		verifier: verifier,
	}
}

type (
	RefreshV1 struct {
		log      *logger.Logger
		issuer   TokenIssuer
		verifier RefreshTokenVerifier
	}

	RefreshV1Request struct {
		RefreshToken string `json:"refresh_token"`
	}

	RefreshV1Response struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenVerifier interface {
		VerifyRefresh(ctx context.Context, refreshToken string) (string, error)
	}
)

func (h *RefreshV1) Path() string {
	return "POST /v1/auth/refresh"
}

func (h *RefreshV1) Handle(ctx context.Context, req *RefreshV1Request) *httpx.Response {
	if req.RefreshToken == "" {
		h.log.Warn("refresh_token is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "refresh_token is required")
	}

	var (
		err    error
		userID string
	)
	if userID, err = h.verifier.VerifyRefresh(ctx, req.RefreshToken); err != nil {
		h.log.Warn("invalid refresh token")
		return httpx.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token")
	}

	var accessToken string
	if accessToken, err = h.issuer.IssueAccessToken(userID); err != nil {
		h.log.Error("can't issue access token", logger.Error(err))
		return httpx.InternalServerError
	}

	var refreshToken string
	if refreshToken, err = h.issuer.IssueRefreshToken(userID); err != nil {
		h.log.Error("can't issue refresh token", logger.Error(err))
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
