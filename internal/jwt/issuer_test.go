package jwt_test

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/riabininkf/http-auth-example/internal/jwt"
)

func TestIssuer_IssueAccessToken(t *testing.T) {
	issuer := jwt.NewIssuer(
		"test_issuer",
		"test_secret",
		time.Second,
		time.Second,
	)

	accessToken, err := issuer.IssueAccessToken("test_user")
	assert.NoError(t, err)

	parts := strings.Split(accessToken, ".")
	assert.Len(t, parts, 3)

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	assert.NoError(t, err)

	header := gjson.ParseBytes(headerBytes)
	assert.Equal(t, "HS256", header.Get("alg").String())
	assert.Equal(t, "JWT", header.Get("typ").String())

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	payload := gjson.ParseBytes(payloadBytes)
	assert.Equal(t, "test_issuer", payload.Get("iss").String())
	assert.Equal(t, "test_user", payload.Get("sub").String())
	assert.Equal(t, "access_token", payload.Get("typ").String())
}

func TestIssuer_IssueRefreshToken(t *testing.T) {
	issuer := jwt.NewIssuer(
		"test_issuer",
		"test_secret",
		time.Second,
		time.Second,
	)

	accessToken, err := issuer.IssueRefreshToken("test_user")
	assert.NoError(t, err)

	parts := strings.Split(accessToken, ".")
	assert.Len(t, parts, 3)

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	assert.NoError(t, err)

	header := gjson.ParseBytes(headerBytes)
	assert.Equal(t, "HS256", header.Get("alg").String())
	assert.Equal(t, "JWT", header.Get("typ").String())

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	payload := gjson.ParseBytes(payloadBytes)
	assert.Equal(t, "test_issuer", payload.Get("iss").String())
	assert.Equal(t, "test_user", payload.Get("sub").String())
	assert.Equal(t, "refresh_token", payload.Get("typ").String())
}
