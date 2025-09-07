package http

import (
	"net/http"

	"github.com/riabininkf/go-modules/logger"
	"github.com/riabininkf/httpx"

	"github.com/riabininkf/http-auth-example/internal/http/handlers"
)

// NewService creates a new *Service instance
func NewService(
	log *logger.Logger,
	loginV1 *handlers.LoginV1,
	refreshV1 *handlers.RefreshV1,
	registerV1 *handlers.RegisterV1,
	updatePasswordV1 *handlers.UpdatePasswordV1,
) *Service {
	return &Service{
		log:              log,
		loginV1:          loginV1,
		refreshV1:        refreshV1,
		registerV1:       registerV1,
		updatePasswordV1: updatePasswordV1,
	}
}

// Service is a facade for http handlers that represents generic handlers as http.HandlerFunc
type Service struct {
	log              *logger.Logger
	loginV1          *handlers.LoginV1
	refreshV1        *handlers.RefreshV1
	registerV1       *handlers.RegisterV1
	updatePasswordV1 *handlers.UpdatePasswordV1
}

// LoginV1 returns http.HandlerFunc for LoginV1 handler
func (s *Service) LoginV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(
		newErrorLogger(s.log), s.loginV1.Handle)
}

// RefreshV1 returns http.HandlerFunc for RefreshV1 handler
func (s *Service) RefreshV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(newErrorLogger(s.log), s.refreshV1.Handle)
}

// RegisterV1 returns http.HandlerFunc for RegisterV1 handler
func (s *Service) RegisterV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(newErrorLogger(s.log), s.registerV1.Handle)
}

// UpdatePasswordV1 returns http.HandlerFunc for UpdatePasswordV1 handler
func (s *Service) UpdatePasswordV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(newErrorLogger(s.log), s.updatePasswordV1.Handle)
}
