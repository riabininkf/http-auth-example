package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riabininkf/go-modules/logger"
)

var (
	ErrTokenMissing = errors.New("jwt token is missing")
)

func NewAuthenticator(
	secret string,
	log *logger.Logger,
	parser JwtParser,
	noAuthRoutes []string,
) *Authenticator {
	methods := make(map[string]struct{}, len(noAuthRoutes))
	for _, method := range noAuthRoutes {
		methods[method] = struct{}{}
	}

	return &Authenticator{
		secret:       secret,
		log:          log,
		parser:       parser,
		noAuthRoutes: methods,
	}
}

type (
	Authenticator struct {
		secret string
		log    *logger.Logger
		parser JwtParser

		noAuthRoutes map[string]struct{}
	}

	JwtParser interface {
		ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	}
)

func (a *Authenticator) Authenticate(_ context.Context, req *http.Request) (string, error) {
	var (
		err   error
		token string
	)
	if token, err = a.parseToken(req.Header.Get("Authorization")); err != nil && a.isAuthRequired(req) {
		return "", err
	}

	return token, nil
}

func (a *Authenticator) parseToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrTokenMissing
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))

	var (
		err         error
		claims      jwt.RegisteredClaims
		parsedToken *jwt.Token
	)
	if parsedToken, err = a.parser.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(a.secret), nil
	}); err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Subject == "" {
		return "", errors.New("subject is missing")
	}

	return claims.Subject, nil
}

func (a *Authenticator) isAuthRequired(req *http.Request) bool {
	_, ok := a.noAuthRoutes[fmt.Sprintf("%s %s", req.Method, req.URL.Path)]
	return !ok
}
