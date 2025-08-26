package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain/auth"
)

func NewUpdatePasswordV1(
	log *logger.Logger,
	userProvider UserByIdProvider,
	passwordUpdater PasswordUpdater,
) *UpdatePasswordV1 {
	return &UpdatePasswordV1{
		log:             log,
		userProvider:    userProvider,
		passwordUpdater: passwordUpdater,
	}
}

type (
	UpdatePasswordV1 struct {
		log             *logger.Logger
		userProvider    UserByIdProvider
		passwordUpdater PasswordUpdater
	}

	UpdatePasswordV1Request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	UserByIdProvider interface {
		GetByID(ctx context.Context, userID string) (auth.User, error)
	}

	PasswordUpdater interface {
		UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
	}
)

func (h *UpdatePasswordV1) Path() string {
	return "POST /v1/user/password"
}

func (h *UpdatePasswordV1) Handle(ctx context.Context, req *UpdatePasswordV1Request) *httpx.Response {
	if req.OldPassword == "" {
		h.log.Warn("old password is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "old_password is required")
	}

	if req.NewPassword == "" {
		h.log.Warn("new password is missing")
		return httpx.NewErrorResponse(http.StatusBadRequest, "new_password is required")
	}

	if req.OldPassword == req.NewPassword {
		h.log.Warn("old and new passwords are the same")
		return httpx.NewJsonResponse(httpx.WithStatus(http.StatusOK))
	}

	// validation is skipped for simplicity

	var (
		ok     bool
		userID string
	)
	if userID, ok = httpx.GetUserID(ctx); !ok {
		h.log.Warn("user id is missing")
		return httpx.BadRequest
	}

	var (
		err  error
		user auth.User
	)
	if user, err = h.userProvider.GetByID(ctx, userID); err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			h.log.Warn("user not found")
			return httpx.NotFound
		}

		h.log.Error("can't get user by id", logger.Error(err))
		return httpx.InternalServerError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword()), []byte(req.OldPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			h.log.Warn("invalid password")
			return httpx.NewErrorResponse(http.StatusBadRequest, "invalid old password")
		}

		h.log.Error("can't compare passwords", logger.Error(err))
		return httpx.InternalServerError
	}

	var hashedPassword []byte
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost); err != nil {
		h.log.Error("can't generate password hash", logger.Error(err))
		return httpx.InternalServerError
	}

	if err = h.passwordUpdater.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		h.log.Error("can't update password", logger.Error(err))
		return httpx.InternalServerError
	}

	return httpx.NewJsonResponse(httpx.WithStatus(http.StatusOK))
}
