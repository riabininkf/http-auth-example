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

func TestUpdatePasswordV1(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		statusCode, resp := sendUpdatePasswordV1Request(t, "", nil)

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "unauthorized", resp.Get("error.message").String())
	})

	t.Run("old password is missing", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		registerUserV1(t, email, password)
		accessToken := loginUserV1(t, email, password)

		statusCode, resp := sendUpdatePasswordV1Request(t, accessToken, nil)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "old_password is required", resp.Get("error.message").String())
	})

	t.Run("new password is missing", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		registerUserV1(t, email, password)
		accessToken := loginUserV1(t, email, password)

		statusCode, resp := sendUpdatePasswordV1Request(t, accessToken, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"old_password":"%s"}`, password)),
		))

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "new_password is required", resp.Get("error.message").String())
	})

	t.Run("old and new passwords are equal", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		registerUserV1(t, email, password)
		accessToken := loginUserV1(t, email, password)

		statusCode, _ := sendUpdatePasswordV1Request(t, accessToken, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"old_password":"%s", "new_password":"%s"}`, password, password)),
		))

		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("invalid old password", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		registerUserV1(t, email, password)
		accessToken := loginUserV1(t, email, password)

		statusCode, resp := sendUpdatePasswordV1Request(t, accessToken, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"old_password":"%s", "new_password":"%s"}`, gofakeit.Name(), gofakeit.Name())),
		))

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, "invalid old password", resp.Get("error.message").String())
	})

	t.Run("positive case", func(t *testing.T) {
		email, password := gofakeit.Email(), gofakeit.Name()

		registerUserV1(t, email, password)
		accessToken := loginUserV1(t, email, password)

		newPassword := gofakeit.Name()
		statusCode, _ := sendUpdatePasswordV1Request(t, accessToken, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"old_password":"%s", "new_password":"%s"}`, password, newPassword)),
		))

		assert.Equal(t, http.StatusOK, statusCode)

		// try to log in with the old password

		var resp gjson.Result
		statusCode, resp = sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, gofakeit.Name())),
		))

		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, "invalid email or password", resp.Get("error.message").String())

		// try to log in with the new password

		statusCode, resp = sendLoginV1Request(t, bytes.NewReader(
			[]byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, newPassword)),
		))

		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func sendUpdatePasswordV1Request(t *testing.T, accessToken string, body io.Reader) (int, gjson.Result) {
	return sendHttpRequest(t, http.MethodPost, "http://localhost:8080/v1/user/password", body, accessToken)
}
