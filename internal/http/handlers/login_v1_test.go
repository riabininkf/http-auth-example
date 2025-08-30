package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/riabininkf/httpx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/riabininkf/http-auth-example/internal/domain"
	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/handlers/mocks"
)

func TestLoginV1_Handle(t *testing.T) {
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

	generateRequest := func() *handlers.LoginV1Request {
		return &handlers.LoginV1Request{Email: gofakeit.Email(), Password: gofakeit.Name()}
	}

	testCases := []struct {
		name                string
		req                 func() *handlers.LoginV1Request
		onGetByEmail        func(req *handlers.LoginV1Request) (domain.User, error)
		onIssueAccessToken  func() (string, error)
		onIssueRefreshToken func() (string, error)
		onSaveRefreshToken  func() error
		expResp             *httpx.Response
	}{
		{
			name:    "email is missing",
			req:     func() *handlers.LoginV1Request { return &handlers.LoginV1Request{} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "email is required"),
		},
		{
			name:    "password is missing",
			req:     func() *handlers.LoginV1Request { return &handlers.LoginV1Request{Email: gofakeit.Email()} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "password is required"),
		},
		{
			name: "user not found",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			expResp: httpx.NewErrorResponse(http.StatusUnauthorized, "invalid email or password"),
		},
		{
			name: "can't get user by email",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return nil, errors.New("test error")
			},
			expResp: httpx.InternalServerError,
		},
		{
			name: "invalid password",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return domain.NewUser(uuid.NewString(), req.Email, generatePasswordHash(t, gofakeit.Name())), nil
			},
			expResp: httpx.NewErrorResponse(http.StatusUnauthorized, "invalid email or password"),
		},
		{
			name: "can't issue access token",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return domain.NewUser(uuid.NewString(), req.Email, generatePasswordHash(t, req.Password)), nil
			},
			onIssueAccessToken: func() (string, error) { return "", errors.New("test error") },
			expResp:            httpx.InternalServerError,
		},
		{
			name: "can't issue refresh token",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return domain.NewUser(uuid.NewString(), req.Email, generatePasswordHash(t, req.Password)), nil
			},
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "", errors.New("test error") },
			expResp:             httpx.InternalServerError,
		},
		{
			name: "can't save refresh token",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return domain.NewUser(uuid.NewString(), req.Email, generatePasswordHash(t, req.Password)), nil
			},
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "refresh_token", nil },
			onSaveRefreshToken:  func() error { return errors.New("test error") },
			expResp:             httpx.InternalServerError,
		},
		{
			name: "positive case",
			req:  generateRequest,
			onGetByEmail: func(req *handlers.LoginV1Request) (domain.User, error) {
				return domain.NewUser("user_id", req.Email, generatePasswordHash(t, req.Password)), nil
			},
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "refresh_token", nil },
			onSaveRefreshToken:  func() error { return nil },
			expResp: httpx.NewJsonResponse(
				httpx.WithStatus(http.StatusOK),
				httpx.WithBody(&handlers.LoginV1Response{
					UserID:       "user_id",
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				}),
			),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			req := testCase.req()

			userProvider := mocks.NewUserByEmailProvider(t)

			var user domain.User
			if testCase.onGetByEmail != nil {
				var err error
				user, err = testCase.onGetByEmail(req)

				userProvider.On("GetByEmail", ctx, req.Email).
					Return(user, err)
			}

			tokenIssuer := mocks.NewTokenIssuer(t)
			if testCase.onIssueAccessToken != nil {
				tokenIssuer.On("IssueAccessToken", user.ID()).Return(testCase.onIssueAccessToken())
			}

			var refreshToken string
			if testCase.onIssueRefreshToken != nil {
				var err error
				refreshToken, err = testCase.onIssueRefreshToken()

				tokenIssuer.On("IssueRefreshToken", user.ID()).Return(refreshToken, err)
			}

			jwtStorage := mocks.NewJwtStorage(t)
			if testCase.onSaveRefreshToken != nil {
				jwtStorage.On("Save", ctx, refreshToken).Return(testCase.onSaveRefreshToken())
			}

			handler := handlers.NewLoginV1(
				zap.NewNop(),
				tokenIssuer,
				jwtStorage,
				userProvider,
			)

			assert.Equal(t, testCase.expResp, handler.Handle(ctx, req))
		})
	}
}
