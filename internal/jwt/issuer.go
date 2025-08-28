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

type (
	Issuer struct {
		issuer          string
		secret          string
		accessTokenTTL  time.Duration
		refreshTokenTTL time.Duration
	}

	claimsWithType struct {
		jwt.RegisteredClaims
		Type string `json:"typ"`
	}
)

func (i *Issuer) IssueAccessToken(userID string) (string, error) {
	return i.issueToken(userID, i.accessTokenTTL, tokenTypeAccessToken)
}

func (i *Issuer) IssueRefreshToken(userID string) (string, error) {
	return i.issueToken(userID, i.refreshTokenTTL, tokenTypeRefreshToken)
}

func (i *Issuer) issueToken(userID string, ttl time.Duration, tokenType string) (string, error) {
	now := time.Now()

	// asymmetric key is better for most cases since only the auth server should know the secret,
	// but symmetric key is used here for simplicity as this is just an example
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsWithType{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		Type: tokenType,
	}).SignedString([]byte(i.secret))
}
