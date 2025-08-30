package middleware

import (
	"context"
	"net/http"

	"github.com/riabininkf/go-modules/logger"
	"github.com/riabininkf/httpx"
)

type Authenticator interface {
	Authenticate(ctx context.Context, req *http.Request) (string, error)
}

func Auth(log *logger.Logger, verifier Authenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			var (
				err    error
				userID string
			)
			if userID, err = verifier.Authenticate(req.Context(), req); err != nil {
				log.Warn("user is not authenticated", logger.Error(err))

				if err = httpx.WriteJsonResponse(httpx.Unauthorized, writer); err != nil {
					log.Error("can't write error response", logger.Error(err))
				}

				return
			}

			next.ServeHTTP(writer, req.WithContext(
				httpx.ContextWithUserID(req.Context(), userID),
			))
		})
	}
}
