package http_transport

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func New(service serviceInterface, logs *zap.Logger) (*Server, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	s := new(Server)
	s.cfg = cfg

	s.service = service
	if service == nil {
		return nil, errors.New("param `service` is required")
	}

	s.log = logs
	if logs == nil {
		return nil, errors.New("param `logs` is required")
	}

	s.a = &http.Server{
		Addr:    cfg.Host,
		Handler: chi.NewRouter(),
	}

	return s.setupHandlers(), nil
}

func (s *Server) setupHandlers() *Server {

	s.a.Handler.(*chi.Mux).Get("/api/ping", s.pingHandler)

	s.a.Handler.(*chi.Mux).Route("/api/v1/", func(r chi.Router) {
		r.Get("/access", s.accessHandler)
		r.Post("/refresh", s.refreshHandler)
	})

	return s
}
