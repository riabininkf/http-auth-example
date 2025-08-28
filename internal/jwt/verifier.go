package jwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

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
	Verifier struct {
		secret string
		parser Parser
	}

	Parser interface {
		ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	}
)

func (v *Verifier) VerifyAccess(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeAccessToken)
}

func (v *Verifier) VerifyRefresh(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeRefreshToken)
}

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
