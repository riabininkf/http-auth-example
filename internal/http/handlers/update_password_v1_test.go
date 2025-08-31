package handlers_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/riabininkf/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain"
	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/handlers/mocks"
)

func TestNewUpdatePasswordV1(t *testing.T) {
	generatePasswordHash := func(t *testing.T, password string) string {
		t.Helper()
		var (
			err            error
			bcryptPassword []byte
		)
		if bcryptPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
			t.Fatal(err)
		}

		return string(bcryptPassword)
	}

	generateRequest := func() *handlers.UpdatePasswordV1Request {
		return &handlers.UpdatePasswordV1Request{
			OldPassword: "old_password",
			NewPassword: gofakeit.Name(),
		}
	}

	testCases := []struct {
		name             string
		req              func() *handlers.UpdatePasswordV1Request
		userID           string
		onGetUserByID    func() (domain.User, error)
		onUpdatePassword func() error
		expResp          *httpx.Response
	}{
		{
			name:    "old password is missing",
			req:     func() *handlers.UpdatePasswordV1Request { return &handlers.UpdatePasswordV1Request{} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "old_password is required"),
		},
		{
			name: "new password is missing",
			req: func() *handlers.UpdatePasswordV1Request {
				return &handlers.UpdatePasswordV1Request{
					OldPassword: gofakeit.Name(),
				}
			},
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "new_password is required"),
		},
		{
			name: "old and new passwords are the same",
			req: func() *handlers.UpdatePasswordV1Request {
				return &handlers.UpdatePasswordV1Request{
					OldPassword: "test_password",
					NewPassword: "test_password",
				}
			},
			expResp: httpx.NewJsonResponse(httpx.WithStatus(http.StatusOK)),
		},
		{
			name:    "user id is missing",
			req:     generateRequest,
			expResp: httpx.BadRequest,
		},
		{
			name:   "user not found",
			req:    generateRequest,
			userID: "user_id",
			onGetUserByID: func() (domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			expResp: httpx.NotFound,
		},
		{
			name:   "can't get user by id",
			req:    generateRequest,
			userID: "user_id",
			onGetUserByID: func() (domain.User, error) {
				return nil, assert.AnError
			},
			expResp: httpx.InternalServerError,
		},
		{
			name:   "invalid password",
			req:    generateRequest,
			userID: "user_id",
			onGetUserByID: func() (domain.User, error) {
				return domain.NewUser(uuid.NewString(), gofakeit.Email(), generatePasswordHash(t, "incorrect")), nil
			},
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "invalid old password"),
		},
		{
			name:   "can't update password",
			req:    generateRequest,
			userID: "user_id",
			onGetUserByID: func() (domain.User, error) {
				return domain.NewUser(uuid.NewString(), gofakeit.Email(), generatePasswordHash(t, "old_password")), nil
			},
			onUpdatePassword: func() error { return assert.AnError },
			expResp:          httpx.InternalServerError,
		},
		{
			name:   "positive case",
			req:    generateRequest,
			userID: "user_id",
			onGetUserByID: func() (domain.User, error) {
				return domain.NewUser(uuid.NewString(), gofakeit.Email(), generatePasswordHash(t, "old_password")), nil
			},
			onUpdatePassword: func() error { return nil },
			expResp:          httpx.NewJsonResponse(httpx.WithStatus(http.StatusOK)),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := testCase.req()

			ctx := t.Context()
			if testCase.userID != "" {
				ctx = httpx.ContextWithUserID(ctx, testCase.userID)
			}

			userProvider := mocks.NewUserByIdProvider(t)
			if testCase.onGetUserByID != nil {
				userProvider.On("GetByID", ctx, testCase.userID).Return(testCase.onGetUserByID())
			}

			passwordUpdater := mocks.NewPasswordUpdater(t)
			if testCase.onUpdatePassword != nil {
				passwordUpdater.On("UpdatePassword", ctx, testCase.userID, mock.AnythingOfType("string")).
					Return(testCase.onUpdatePassword())
			}

			handler := handlers.NewUpdatePasswordV1(
				zap.NewNop(),
				userProvider,
				passwordUpdater,
			)

			assert.Equal(t, testCase.expResp, handler.Handle(ctx, req))
		})
	}
}
