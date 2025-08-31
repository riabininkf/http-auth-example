package jwt_test

import (
	"context"
	"reflect"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/riabininkf/http-auth-example/internal/jwt"
	"github.com/riabininkf/http-auth-example/internal/jwt/mocks"
)

func setClaims(t *testing.T, c gojwt.Claims, tokenType, subject string) {
	t.Helper()
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		t.Fatalf("claims is not a pointer: %T", c)
	}
	s := v.Elem()
	if !s.IsValid() {
		t.Fatalf("claims value is not valid")
	}

	typeField := s.FieldByName("Type")
	if !typeField.IsValid() || !typeField.CanSet() {
		t.Fatalf("claims does not have settable field 'Type'")
	}
	typeField.SetString(tokenType)

	subField := s.FieldByName("Subject")
	if !subField.IsValid() || !subField.CanSet() {
		t.Fatalf("claims does not have settable field 'Subject'")
	}
	subField.SetString(subject)
}

func TestVerifier_VerifyAccess(t *testing.T) {
	t.Run("parser returns an error", func(t *testing.T) {
		parser := mocks.NewParser(t)

		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return((*gojwt.Token)(nil), assert.AnError)

		subject, err := jwt.NewVerifier("secret", parser).
			VerifyAccess(context.Background(), "token")

		assert.Empty(t, subject)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("invalid signing method", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
			_, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodNone})
			return nil, err
		})

		subject, err := jwt.NewVerifier("secret", parser).
			VerifyAccess(context.Background(), "token")

		assert.Empty(t, subject)
		assert.ErrorIs(t, err, gojwt.ErrSignatureInvalid)
	})

	t.Run("parsed token is invalid", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
			if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
				return nil, err
			}

			return &gojwt.Token{Valid: false, Method: gojwt.SigningMethodHS256}, nil
		})

		subject, err := jwt.NewVerifier("secret", parser).
			VerifyAccess(context.Background(), "token")

		assert.Empty(t, subject)
		assert.EqualError(t, err, "invalid token")
	})

	t.Run("claims type mismatch", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
			if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
				return nil, err
			}

			setClaims(t, _claims, "refresh_token", "user_id")
			return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
		})

		subject, err := jwt.NewVerifier("secret", parser).
			VerifyAccess(context.Background(), "token")

		assert.Empty(t, subject)
		assert.ErrorIs(t, err, gojwt.ErrTokenInvalidClaims)
	})

	t.Run("subject is empty", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
			if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
				return nil, err
			}
			setClaims(t, _claims, "access_token", "")
			return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
		})

		verifier := jwt.NewVerifier("secret", parser)

		subject, err := verifier.VerifyAccess(context.Background(), "token")
		assert.Empty(t, subject)
		assert.ErrorIs(t, err, gojwt.ErrTokenInvalidClaims)
	})

	t.Run("positive case", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims",
			"token",
			mock.AnythingOfType("*jwt.claimsWithType"),
			mock.AnythingOfType("jwt.Keyfunc"),
		).Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
			if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
				return nil, err
			}
			setClaims(t, _claims, "access_token", "user_123")
			return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
		})

		verifier := jwt.NewVerifier("secret", parser)

		subject, err := verifier.VerifyAccess(context.Background(), "token")
		assert.NoError(t, err)
		assert.Equal(t, "user_123", subject)
	})
}

func TestVerifier_VerifyRefresh(t *testing.T) {
	t.Run("claims type mismatch", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims", mock.Anything, mock.Anything, mock.Anything).
			Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
				if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
					return nil, err
				}
				setClaims(t, _claims, "access_token", "user_id")
				return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
			})

		verifier := jwt.NewVerifier("secret", parser)

		subject, err := verifier.VerifyRefresh(context.Background(), "token")
		assert.Empty(t, subject)
		assert.ErrorIs(t, err, gojwt.ErrTokenInvalidClaims)
	})

	t.Run("subject is empty", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims", mock.Anything, mock.Anything, mock.Anything).
			Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
				if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
					return nil, err
				}
				setClaims(t, _claims, "refresh_token", "")
				return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
			})

		verifier := jwt.NewVerifier("secret", parser)

		subject, err := verifier.VerifyRefresh(context.Background(), "token")
		assert.Empty(t, subject)
		assert.ErrorIs(t, err, gojwt.ErrTokenInvalidClaims)
	})

	t.Run("positive case", func(t *testing.T) {
		parser := mocks.NewParser(t)
		parser.On("ParseWithClaims", mock.Anything, mock.Anything, mock.Anything).
			Return(func(_ string, _claims gojwt.Claims, keyFunc gojwt.Keyfunc) (*gojwt.Token, error) {
				if _, err := keyFunc(&gojwt.Token{Method: gojwt.SigningMethodHS256}); err != nil {
					return nil, err
				}
				setClaims(t, _claims, "refresh_token", "user_456")
				return &gojwt.Token{Valid: true, Method: gojwt.SigningMethodHS256}, nil
			})

		verifier := jwt.NewVerifier("secret", parser)

		subject, err := verifier.VerifyRefresh(context.Background(), "token")
		assert.NoError(t, err)
		assert.Equal(t, "user_456", subject)
	})
}
