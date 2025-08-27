package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewIssuer(
	issuer string,
	secret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Issuer {
	return &Issuer{
		issuer:          issuer,
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

const (
	tokenTypeAccessToken  = "access_token"
	tokenTypeRefreshToken = "refresh_token"
)

type Issuer struct {
	issuer          string
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (i *Issuer) IssueAccessToken(userID string) (string, error) {
	return i.issueToken(userID, i.accessTokenTTL, tokenTypeAccessToken)
}

func (i *Issuer) IssueRefreshToken(userID string) (string, error) {
	return i.issueToken(userID, i.refreshTokenTTL, tokenTypeRefreshToken)
}

func (i *Issuer) issueToken(userID string, ttl time.Duration, tokenType string) (string, error) {
	now := time.Now()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": i.issuer,
		"sub": userID,
		"iat": jwt.NewNumericDate(now).Unix(),
		"exp": jwt.NewNumericDate(now.Add(ttl)).Unix(),
		"typ": tokenType,
	}).SignedString([]byte(i.secret))
}
