package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/riabininkf/go-modules/httpx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/handlers/mocks"
)

func TestRefreshV1_Handle(t *testing.T) {
	generateRequest := func() *handlers.RefreshV1Request {
		return &handlers.RefreshV1Request{RefreshToken: gofakeit.Name()}
	}

	testCases := []struct {
		name                string
		req                 func() *handlers.RefreshV1Request
		onPop               func() error
		onVerifyRefresh     func() (string, error)
		onIssueAccessToken  func() (string, error)
		onIssueRefreshToken func() (string, error)
		onSaveRefreshToken  func() error
		expResp             *httpx.Response
	}{
		{
			name:    "refresh token is missing",
			req:     func() *handlers.RefreshV1Request { return &handlers.RefreshV1Request{} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "refresh_token is required"),
		},
		{
			name:    "failed to pop refresh token from the storage",
			req:     generateRequest,
			onPop:   func() error { return errors.New("test error") },
			expResp: httpx.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token"),
		},
		{
			name:            "failed to verify refresh token",
			req:             generateRequest,
			onPop:           func() error { return nil },
			onVerifyRefresh: func() (string, error) { return "", errors.New("test error") },
			expResp:         httpx.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token"),
		},
		{
			name:               "can't issue access token",
			req:                generateRequest,
			onPop:              func() error { return nil },
			onVerifyRefresh:    func() (string, error) { return gofakeit.Name(), nil },
			onIssueAccessToken: func() (string, error) { return "", errors.New("test error") },
			expResp:            httpx.InternalServerError,
		},
		{
			name:                "can't issue refresh token",
			req:                 generateRequest,
			onPop:               func() error { return nil },
			onVerifyRefresh:     func() (string, error) { return gofakeit.Name(), nil },
			onIssueAccessToken:  func() (string, error) { return gofakeit.Name(), nil },
			onIssueRefreshToken: func() (string, error) { return "", errors.New("test error") },
			expResp:             httpx.InternalServerError,
		},
		{
			name:                "can't save refresh token",
			req:                 generateRequest,
			onPop:               func() error { return nil },
			onVerifyRefresh:     func() (string, error) { return gofakeit.Name(), nil },
			onIssueAccessToken:  func() (string, error) { return gofakeit.Name(), nil },
			onIssueRefreshToken: func() (string, error) { return gofakeit.Name(), nil },
			onSaveRefreshToken:  func() error { return errors.New("test error") },
			expResp:             httpx.InternalServerError,
		},
		{
			name:                "positive case",
			req:                 generateRequest,
			onPop:               func() error { return nil },
			onVerifyRefresh:     func() (string, error) { return "user_id", nil },
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "refresh_token", nil },
			onSaveRefreshToken:  func() error { return nil },
			expResp: httpx.NewJsonResponse(
				httpx.WithStatus(http.StatusOK),
				httpx.WithBody(&handlers.RefreshV1Response{
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

			jwtStorage := mocks.NewJwtStorage(t)
			if testCase.onPop != nil {
				jwtStorage.On("Pop", ctx, req.RefreshToken).Return(testCase.onPop())
			}

			refreshVerifier := mocks.NewRefreshTokenVerifier(t)

			var userID string
			if testCase.onVerifyRefresh != nil {
				var err error
				userID, err = testCase.onVerifyRefresh()

				refreshVerifier.On("VerifyRefresh", ctx, req.RefreshToken).Return(userID, err)
			}

			issuer := mocks.NewTokenIssuer(t)

			var accessToken string
			if testCase.onIssueAccessToken != nil {
				var err error
				accessToken, err = testCase.onIssueAccessToken()

				issuer.On("IssueAccessToken", userID).Return(accessToken, err)
			}

			var refreshToken string
			if testCase.onIssueRefreshToken != nil {
				var err error
				refreshToken, err = testCase.onIssueRefreshToken()

				issuer.On("IssueRefreshToken", userID).Return(refreshToken, err)
			}

			if testCase.onSaveRefreshToken != nil {
				jwtStorage.On("Save", ctx, refreshToken).Return(testCase.onSaveRefreshToken())
			}

			handler := handlers.NewRefreshV1(
				zap.NewNop(),
				issuer,
				jwtStorage,
				refreshVerifier,
			)

			assert.Equal(t, testCase.expResp, handler.Handle(ctx, req))
		})
	}
}
