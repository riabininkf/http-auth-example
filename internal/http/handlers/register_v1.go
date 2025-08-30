package handlers

//go:generate mockery --name UserRegistrar --output ./mocks --outpkg mocks --filename user_registrar.go --structname UserRegistrar

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/riabininkf/go-modules/logger"
	"github.com/riabininkf/httpx"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain"
)

func NewRegisterV1(
	log *logger.Logger,
	issuer TokenIssuer,
	jwtStorage JwtStorage,
	registrar UserRegistrar,
) *RegisterV1 {
	return &RegisterV1{
		log:        log,
		issuer:     issuer,
		jwtStorage: jwtStorage,
		registrar:  registrar,
	}
}

type (
	RegisterV1 struct {
		log        *logger.Logger
		issuer     TokenIssuer
		jwtStorage JwtStorage
		registrar  UserRegistrar
	}

	RegisterV1Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterV1Response struct {
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UserRegistrar interface {
		Save(ctx context.Context, user domain.User) error
	}
)

func (h *RegisterV1) Handle(ctx context.Context, req *RegisterV1Request) *httpx.Response {
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

	user := domain.NewUser(
		uuid.NewString(),
		req.Email,
		string(hashedPassword),
	)

	if err = h.registrar.Save(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailBusy) {
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

	if err = h.jwtStorage.Save(ctx, refreshToken); err != nil {
		h.log.Error("can't save refresh token", logger.Error(err))
		return httpx.InternalServerError
	}

	return httpx.NewJsonResponse(
		httpx.WithStatus(http.StatusCreated),
		httpx.WithBody(&RegisterV1Response{
			UserID:       user.ID(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}),
	)
}
