package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewIssuer initializes a new Issuer instance with the specified parameters for token generation and expiration settings.
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

// tokenTypeAccessToken represents the type string for an access token.
// tokenTypeRefreshToken represents the type string for a refresh token.
const (
	tokenTypeAccessToken  = "access_token"
	tokenTypeRefreshToken = "refresh_token"
)

type (
	// Issuer represents the structure for storing token issuer configurations and TTLs for access and refresh tokens.
	// claimsWithType extends jwt.RegisteredClaims to include a Type field for specifying claim types.
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

// IssueAccessToken generates a signed access token for the specified user ID with a preset expiration time.
func (i *Issuer) IssueAccessToken(userID string) (string, error) {
	return i.issueToken(userID, i.accessTokenTTL, tokenTypeAccessToken)
}

// IssueRefreshToken generates a new refresh token for the given user ID using the configured TTL and secret.
func (i *Issuer) IssueRefreshToken(userID string) (string, error) {
	return i.issueToken(userID, i.refreshTokenTTL, tokenTypeRefreshToken)
}

// issueToken generates a signed JWT token with a specified TTL and type for the given user ID, using the Issuer's secret key.
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
