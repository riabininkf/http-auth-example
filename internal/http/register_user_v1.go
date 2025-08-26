package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain/auth"
)

func NewRegisterUserV1(
	log *logger.Logger,
	issuer TokenIssuer,
	registrar UserRegistrar,
) *RegisterUserV1 {
	return &RegisterUserV1{
		log:       log,
		issuer:    issuer,
		registrar: registrar,
	}
}

type (
	RegisterUserV1 struct {
		log       *logger.Logger
		issuer    TokenIssuer
		registrar UserRegistrar
	}

	RegisterUserV1Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterUserV1Response struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UserRegistrar interface {
		Save(ctx context.Context, user auth.User) error
	}
)

func (h *RegisterUserV1) Path() string {
	return "POST /v1/auth/register"
}

func (h *RegisterUserV1) Handle(ctx context.Context, req *RegisterUserV1Request) *httpx.Response {
	if req.Email == "" {
		h.log.Warn("email is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "email is required")
	}

	if req.Password == "" {
		h.log.Warn("password is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "password is required")
	}

	// validation is skipped for simplicity

	var (
		err            error
		hashedPassword []byte
	)
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		h.log.Error("can't generate password hash", logger.Error(err))
		return httpx.InternalServerError
	}

	user := auth.NewUser(
		uuid.NewString(),
		req.Email,
		string(hashedPassword),
	)

	if err = h.registrar.Save(ctx, user); err != nil {
		if errors.Is(err, auth.ErrEmailBusy) {
			h.log.Warn("user already exists")
			return httpx.NewErrorResponse(http.StatusBadRequest, "user already exists")
		}

		h.log.Error("can't save user", logger.Error(err))
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
	}

	return httpx.NewJsonResponse(
		httpx.WithStatus(http.StatusCreated),
		httpx.WithBody(&RegisterUserV1Response{
			UserID:       user.ID(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}),
	)
}
