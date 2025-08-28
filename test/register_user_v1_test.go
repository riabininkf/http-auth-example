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

	handlers "github.com/riabininkf/http-auth-example/internal/http"
)

func TestRegisterUserV1(t *testing.T) {
	t.Run("email is missing", func(t *testing.T) {
		statusCode, resp := sendRegistrationV1Request(t, nil)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "email is required", resp.Get("error.message").String())
	})

	t.Run("password is missing", func(t *testing.T) {
		statusCode, resp := sendRegistrationV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s"}`, gofakeit.Email())),
		))

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "password is required", resp.Get("error.message").String())
	})

	t.Run("positive case", func(t *testing.T) {
		statusCode, resp := sendRegistrationV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s", "password":"1234"}`, gofakeit.Email())),
		))

		assert.Equal(t, http.StatusCreated, statusCode)

		assert.True(t, resp.Get("user_id").Exists(), "user_id is missing")
		assert.True(t, resp.Get("access_token").Exists(), "access_token is missing")
		assert.True(t, resp.Get("refresh_token").Exists(), "refresh_token is missing")
	})

	t.Run("user already exists", func(t *testing.T) {
		email := gofakeit.Email()

		statusCode, resp := sendRegistrationV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s", "password":"1234"}`, email)),
		))

		assert.Equal(t, http.StatusCreated, statusCode)

		assert.True(t, resp.Get("user_id").Exists(), "user_id is missing")
		assert.True(t, resp.Get("access_token").Exists(), "access_token is missing")
		assert.True(t, resp.Get("refresh_token").Exists(), "refresh_token is missing")

		statusCode, resp = sendRegistrationV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s", "password":"1234"}`, email)),
		))

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "user already exists", resp.Get("error.message").String())
	})
}

func sendRegistrationV1Request(t *testing.T, body io.Reader) (int, gjson.Result) {
	return sendHttpRequest(t, http.MethodPost, "http://localhost:8080/v1/auth/register", body, "")
}

func registerUserV1(t *testing.T, email string, password string) *handlers.RegisterUserV1Response {
	statusCode, resp := sendRegistrationV1Request(t, bytes.NewReader(
		[]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)),
	))

	assert.Equal(t, http.StatusCreated, statusCode)

	assert.NotEmpty(t, resp.Get("user_id").String(), "user_id is missing")
	assert.NotEmpty(t, resp.Get("access_token").String(), "access_token is missing")
	assert.NotEmpty(t, resp.Get("refresh_token").String(), "refresh_token is missing")

	return &handlers.RegisterUserV1Response{
		UserID:       resp.Get("user_id").String(),
		AccessToken:  resp.Get("access_token").String(),
		RefreshToken: resp.Get("refresh_token").String(),
	}
}
