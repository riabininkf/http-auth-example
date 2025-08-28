package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestRefreshV1(t *testing.T) {
	t.Run("refresh token is missing", func(t *testing.T) {
		statusCode, resp := sendRefreshV1Request(t, nil)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "refresh_token is required", resp.Get("error.message").String())
	})

	t.Run("non-existent  refresh token", func(t *testing.T) {
		nonExistentRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYzNzI4ODksImlhdCI6MTc1NjM2OTI4OSwiaXNzIjoiYXV0aC1zZXJ2aWNlIiwic3ViIjoiYmYyYjZhNmEtMmZkOC00MzQ0LWE3YzctMThhZmE2ZjJiNDhhIiwidHlwIjoicmVmcmVzaF90b2tlbiJ9.ZZew9QElWO2bLUh0FS51wtw8JBk9K46QJ825NmrqOYk"

		statusCode, resp := sendRefreshV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"refresh_token":"%s"}`, nonExistentRefreshToken)),
		))

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid refresh token", resp.Get("error.message").String())
	})

	t.Run("modified refresh token", func(t *testing.T) {
		registrationResp := registerUserV1(t, gofakeit.Email(), gofakeit.Name())

		// modify the refresh token to provoke a signature error
		parts := strings.Split(registrationResp.RefreshToken, ".")
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		assert.NoError(t, err)

		var claims map[string]any
		assert.NoError(t, json.Unmarshal(payloadBytes, &claims))

		// place a random string instead of userID
		claims["sub"] = gofakeit.Name()

		payloadBytes, err = json.Marshal(claims)
		assert.NoError(t, err)

		parts[1] = base64.RawURLEncoding.EncodeToString(payloadBytes)

		statusCode, resp := sendRefreshV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"refresh_token":"%s"}`, strings.Join(parts, "."))),
		))

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid refresh token", resp.Get("error.message").String())
	})

	t.Run("positive case", func(t *testing.T) {
		registrationResp := registerUserV1(t, gofakeit.Email(), gofakeit.Name())

		statusCode, resp := sendRefreshV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"refresh_token":"%s"}`, registrationResp.RefreshToken)),
		))

		assert.Equal(t, http.StatusOK, statusCode)
		assert.True(t, resp.Get("access_token").Exists(), "access_token is missing")
		assert.True(t, resp.Get("refresh_token").Exists(), "refresh_token is missing")
	})
}

func sendRefreshV1Request(t *testing.T, body io.Reader) (int, gjson.Result) {
	return sendHttpRequest(t, http.MethodPost, "http://localhost:8080/v1/auth/refresh", body, "")
}
