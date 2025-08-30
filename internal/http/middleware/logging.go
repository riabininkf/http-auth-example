package middleware

import (
	"net/http"

	"github.com/riabininkf/go-modules/logger"
)

func Logging(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("incoming http request",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})
	}
}
