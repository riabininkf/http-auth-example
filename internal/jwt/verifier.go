package jwt

//go:generate mockery --name Parser --output ./mocks --outpkg mocks --filename parser.go --structname Parser

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// NewVerifier creates and returns a new instance of Verifier with the specified secret and Parser implementation.
func NewVerifier(
	secret string,
	parser Parser,
) *Verifier {
	return &Verifier{
		secret: secret,
		parser: parser,
	}
}

type (
	// Verifier is a struct that holds a secret and a parser for verifying JWT tokens.
	Verifier struct {
		secret string
		parser Parser
	}

	// Parser defines an interface for parsing JWT tokens with claims and a key function.
	Parser interface {
		ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	}
)

// VerifyAccess validates an access token and returns the subject if the token is valid, or an error if it is invalid.
func (v *Verifier) VerifyAccess(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeAccessToken)
}

// VerifyRefresh validates a given refresh token and returns the subject if valid, or an error otherwise.
func (v *Verifier) VerifyRefresh(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeRefreshToken)
}

// verify validates a token's signature, claims, and type, and returns the subject if valid or an error otherwise.
func (v *Verifier) verify(_ context.Context, token string, tokenType string) (string, error) {
	var (
		err         error
		claims      claimsWithType
		parsedToken *jwt.Token
	)
	if parsedToken, err = v.parser.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(v.secret), nil
	}); err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Type != tokenType {
		return "", jwt.ErrTokenInvalidClaims
	}

	if claims.Subject == "" {
		return "", jwt.ErrTokenInvalidClaims
	}

	return claims.Subject, nil
}
