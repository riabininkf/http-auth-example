package handlers_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/riabininkf/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/riabininkf/http-auth-example/internal/domain"
	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/handlers/mocks"
)

func TestRegisterV1_Handle(t *testing.T) {
	generateRequest := func() *handlers.RegisterV1Request {
		return &handlers.RegisterV1Request{Email: gofakeit.Email(), Password: gofakeit.Name()}
	}

	testCases := []struct {
		name                string
		req                 func() *handlers.RegisterV1Request
		onSaveUser          func() error
		expResp             *httpx.Response
		onIssueAccessToken  func() (string, error)
		onIssueRefreshToken func() (string, error)
		onSaveRefreshToken  func() error
	}{
		{
			name:    "email is missing",
			req:     func() *handlers.RegisterV1Request { return &handlers.RegisterV1Request{} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "email is required"),
		},
		{
			name:    "password is missing",
			req:     func() *handlers.RegisterV1Request { return &handlers.RegisterV1Request{Email: gofakeit.Email()} },
			expResp: httpx.NewErrorResponse(http.StatusBadRequest, "password is required"),
		},
		{
			name:       "user already exists",
			req:        generateRequest,
			onSaveUser: func() error { return domain.ErrEmailBusy },
			expResp:    httpx.NewErrorResponse(http.StatusBadRequest, "user already exists"),
		},
		{
			name:       "can't save user",
			req:        generateRequest,
			onSaveUser: func() error { return assert.AnError },
			expResp:    httpx.InternalServerError,
		},
		{
			name:               "can't issue access token",
			req:                generateRequest,
			onSaveUser:         func() error { return nil },
			onIssueAccessToken: func() (string, error) { return "", assert.AnError },
			expResp:            httpx.InternalServerError,
		},
		{
			name:                "can't issue refresh token",
			req:                 generateRequest,
			onSaveUser:          func() error { return nil },
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "", assert.AnError },
			expResp:             httpx.InternalServerError,
		},
		{
			name:                "can't save refresh token",
			req:                 generateRequest,
			onSaveUser:          func() error { return nil },
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "refresh_token", nil },
			onSaveRefreshToken:  func() error { return assert.AnError },
			expResp:             httpx.InternalServerError,
		},
		{
			name:                "positive case",
			req:                 generateRequest,
			onSaveUser:          func() error { return nil },
			onIssueAccessToken:  func() (string, error) { return "access_token", nil },
			onIssueRefreshToken: func() (string, error) { return "refresh_token", nil },
			onSaveRefreshToken:  func() error { return nil },
			expResp: httpx.NewJsonResponse(
				httpx.WithStatus(http.StatusCreated),
				httpx.WithBody(&handlers.RegisterV1Response{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				}),
			),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := testCase.req()

			registrar := mocks.NewUserRegistrar(t)
			if testCase.onSaveUser != nil {
				registrar.On("Save", t.Context(), mock.AnythingOfType("*domain.user")).Return(testCase.onSaveUser())
			}

			issuer := mocks.NewTokenIssuer(t)
			if testCase.onIssueAccessToken != nil {
				issuer.On("IssueAccessToken", mock.AnythingOfType("string")).Return(testCase.onIssueAccessToken())
			}

			var refreshToken string
			if testCase.onIssueRefreshToken != nil {
				var err error
				refreshToken, err = testCase.onIssueRefreshToken()

				issuer.On("IssueRefreshToken", mock.AnythingOfType("string")).Return(refreshToken, err)
			}

			jwtStorage := mocks.NewJwtStorage(t)
			if testCase.onSaveRefreshToken != nil {
				jwtStorage.On("Save", t.Context(), refreshToken).Return(testCase.onSaveRefreshToken())
			}

			handler := handlers.NewRegisterV1(
				zap.NewNop(),
				issuer,
				jwtStorage,
				registrar,
			)

			resp := handler.Handle(t.Context(), req)
			assert.Equal(t, testCase.expResp.Status(), resp.Status())
			assert.Equal(t, testCase.expResp.Headers(), resp.Headers())

			if registerV1Resp, ok := testCase.expResp.Body().(*handlers.RegisterV1Response); ok {
				// userID is generated inside the handler, so it is impossible to compare it
				assert.Equal(t, registerV1Resp.AccessToken, resp.Body().(*handlers.RegisterV1Response).AccessToken)
				assert.Equal(t, registerV1Resp.RefreshToken, resp.Body().(*handlers.RegisterV1Response).RefreshToken)
			} else {
				assert.Equal(t, testCase.expResp.Body(), resp.Body())
			}
		})
	}
}
