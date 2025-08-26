package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestLoginV1(t *testing.T) {
	t.Run("email is missing", func(t *testing.T) {
		statusCode, resp := sendLoginV1Request(t, nil)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "email is required", resp.Get("error.message").String())
	})

	t.Run("password is missing", func(t *testing.T) {
		statusCode, resp := sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s"}`, gofakeit.Email())),
		))

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "password is required", resp.Get("error.message").String())
	})

	t.Run("invalid email", func(t *testing.T) {
		statusCode, resp := sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, gofakeit.Email(), gofakeit.Name())),
		))

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid email or password", resp.Get("error.message").String())
	})

	t.Run("invalid password", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()
		registerUserV1(t, email, password)

		statusCode, resp := sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, gofakeit.Name())),
		))

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid email or password", resp.Get("error.message").String())
	})

	t.Run("positive case", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		userID := registerUserV1(t, email, password)

		statusCode, resp := sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)),
		))

		assert.Equal(t, http.StatusOK, statusCode)

		assert.Equal(t, userID, resp.Get("user_id").String(), "different user_id")
		assert.True(t, resp.Get("access_token").Exists(), "access_token is missing")
		assert.True(t, resp.Get("refresh_token").Exists(), "refresh_token is missing")
	})
}

func sendLoginV1Request(t *testing.T, body io.Reader) (int, gjson.Result) {
	return sendHttpRequest(t, http.MethodPost, "http://localhost:8080/v1/auth/login", body, "")
}

func loginUserV1(t *testing.T, email string, password string) string {
	statusCode, resp := sendLoginV1Request(t, bytes.NewReader(
		[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)),
	))

	assert.Equal(t, http.StatusOK, statusCode)

	assert.True(t, resp.Get("user_id").Exists(), "user_id is missing")
	assert.True(t, resp.Get("access_token").Exists(), "access_token is missing")
	assert.True(t, resp.Get("refresh_token").Exists(), "refresh_token is missing")

	return resp.Get("access_token").String()
}
