package http

import (
	"net/http"

	"github.com/riabininkf/httpx"

	"github.com/riabininkf/http-auth-example/internal/http/handlers"
	"github.com/riabininkf/http-auth-example/internal/http/log"
)

func NewService(
	log *log.ErrorLogger,
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
	log              *log.ErrorLogger
	loginV1          *handlers.LoginV1
	refreshV1        *handlers.RefreshV1
	registerV1       *handlers.RegisterV1
	updatePasswordV1 *handlers.UpdatePasswordV1
}

func (s *Service) LoginV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(s.log.WithMethod("LoginV1"), s.loginV1.Handle)
}

func (s *Service) RefreshV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(s.log.WithMethod("RefreshV1"), s.refreshV1.Handle)
}

func (s *Service) RegisterV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(s.log.WithMethod("RegisterV1"), s.registerV1.Handle)
}

func (s *Service) UpdatePasswordV1() http.HandlerFunc {
	return httpx.AdaptHandlerFunc(s.log.WithMethod("UpdatePasswordV1"), s.updatePasswordV1.Handle)
}
