package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riabininkf/http-auth-example/internal/jwt"
	"github.com/riabininkf/http-auth-example/internal/jwt/mocks"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	testCases := []struct {
		name           string
		req            func() *http.Request
		noAuthUrls     []string
		onVerifyAccess func() (string, error)
		expUserID      string
		expError       error
	}{
		{
			name: "header is missing (auth required)",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			expError: jwt.ErrTokenMissing,
		},
		{
			name: "header is missing (no auth required)",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			noAuthUrls: []string{"GET /test"},
			expError:   nil,
		},
		{
			name: "bearer is empty (auth required)",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer ")
				return req
			},
			expError: jwt.ErrTokenMissing,
		},
		{
			name: "bearer is empty (no auth required)",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer ")
				return req
			},
			noAuthUrls: []string{"GET /test"},
			expError:   nil,
		},
		{
			name: "verification failed (auth required)",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer test-token")
				return req
			},
			onVerifyAccess: func() (string, error) { return "", assert.AnError },
			expError:       assert.AnError,
		},
		{
			name: "verification failed (no auth required)",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer test-token")
				return req
			},
			onVerifyAccess: func() (string, error) { return "", assert.AnError },
			noAuthUrls:     []string{"GET /test"},
			expError:       nil,
		},
		{
			name: "positive case",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer test-token")
				return req
			},
			onVerifyAccess: func() (string, error) { return "user_id", nil },
			expUserID:      "user_id",
			expError:       nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			verifier := mocks.NewAccessTokenVerifier(t)
			if testCase.onVerifyAccess != nil {
				verifier.On("VerifyAccess", t.Context(), "test-token").Return(testCase.onVerifyAccess())
			}

			authenticator := jwt.NewAuthenticator(verifier, testCase.noAuthUrls)

			userID, err := authenticator.Authenticate(t.Context(), testCase.req())
			assert.Equal(t, testCase.expUserID, userID)
			assert.Equal(t, testCase.expError, err)
		})
	}
}
