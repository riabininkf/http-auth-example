package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func NewTokenVerifier(
	secret string,
	parser JwtParser,
) *TokenVerifier {
	return &TokenVerifier{
		secret: secret,
		parser: parser,
	}
}

type (
	TokenVerifier struct {
		secret string
		parser JwtParser
	}

	JwtParser interface {
		ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	}
)

func (v *TokenVerifier) VerifyAccess(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeAccessToken)
}

func (v *TokenVerifier) VerifyRefresh(ctx context.Context, token string) (string, error) {
	return v.verify(ctx, token, tokenTypeRefreshToken)
}

func (v *TokenVerifier) verify(_ context.Context, token string, tokenType string) (string, error) {
	var (
		err         error
		claims      jwt.MapClaims
		parsedToken *jwt.Token
	)
	if parsedToken, err = v.parser.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
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

	if parsedType, ok := claims["typ"].(string); !ok || parsedType != tokenType {
		return "", jwt.ErrTokenInvalidClaims
	}

	var (
		ok      bool
		subject string
	)
	if subject, ok = claims["sub"].(string); !ok || subject == "" {
		return "", jwt.ErrTokenInvalidClaims
	}

	return subject, nil
}
