package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewTokenIssuer(
	issuer string,
	secret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *TokenIssuer {
	return &TokenIssuer{
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

type TokenIssuer struct {
	issuer          string
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (i *TokenIssuer) IssueAccessToken(userID string) (string, error) {
	return i.issueToken(userID, i.accessTokenTTL, tokenTypeAccessToken)
}

func (i *TokenIssuer) IssueRefreshToken(userID string) (string, error) {
	return i.issueToken(userID, i.refreshTokenTTL, tokenTypeRefreshToken)
}

func (i *TokenIssuer) issueToken(userID string, ttl time.Duration, tokenType string) (string, error) {
	now := time.Now()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": i.issuer,
		"sub": userID,
		"iat": jwt.NewNumericDate(now).Unix(),
		"exp": jwt.NewNumericDate(now.Add(ttl)).Unix(),
		"typ": tokenType,
	}).SignedString([]byte(i.secret))
}
