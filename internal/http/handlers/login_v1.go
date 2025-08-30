package handlers

//go:generate mockery --name UserByEmailProvider --output ./mocks --outpkg mocks --filename user_by_email_provider.go --structname UserByEmailProvider

import (
	"context"
	"errors"
	"net/http"

	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain"
)

func NewLoginV1(
	log *logger.Logger,
	issuer TokenIssuer,
	jwtStorage JwtStorage,
	userProvider UserByEmailProvider,
) *LoginV1 {
	return &LoginV1{
		log:          log,
		issuer:       issuer,
		jwtStorage:   jwtStorage,
		userProvider: userProvider,
	}
}

type (
	LoginV1 struct {
		log          *logger.Logger
		issuer       TokenIssuer
		jwtStorage   JwtStorage
		userProvider UserByEmailProvider
	}

	LoginV1Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginV1Response struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UserByEmailProvider interface {
		GetByEmail(ctx context.Context, email string) (domain.User, error)
	}
)

func (h *LoginV1) Handle(ctx context.Context, req *LoginV1Request) *httpx.Response {
	if req.Email == "" {
		h.log.Warn("email is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "email is required")
	}

	if req.Password == "" {
		h.log.Warn("password is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "password is required")
	}

	var (
		err  error
		user domain.User
	)
	if user, err = h.userProvider.GetByEmail(ctx, req.Email); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.log.Warn("invalid email")
			return httpx.NewErrorResponse(http.StatusUnauthorized, "invalid email or password")
		}

		h.log.Error("can't get user by email", logger.Error(err))
		return httpx.InternalServerError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword()), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			h.log.Warn("invalid password")
			return httpx.NewErrorResponse(http.StatusUnauthorized, "invalid email or password")
		}

		h.log.Error("can't compare password", logger.Error(err))
		return httpx.InternalServerError
	}

	var accessToken string
	if accessToken, err = h.issuer.IssueAccessToken(user.ID()); err != nil {
		h.log.Error("can't issue access token", logger.Error(err))
		return httpx.InternalServerError
	}

	var refreshToken string
	if refreshToken, err = h.issuer.IssueRefreshToken(user.ID()); err != nil {
		h.log.Error("can't issue refresh token", logger.Error(err))
		return httpx.InternalServerError
	}

	if err = h.jwtStorage.Save(ctx, refreshToken); err != nil {
		h.log.Error("can't save refresh token", logger.Error(err))
		return httpx.InternalServerError
	}

	return httpx.NewJsonResponse(
		httpx.WithStatus(http.StatusOK),
		httpx.WithBody(&LoginV1Response{
			UserID:       user.ID(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}),
	)
}
