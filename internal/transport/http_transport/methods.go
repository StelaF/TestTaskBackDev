package http_transport

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (s *Server) Listen() error {
	var err error

	if (s.cfg.TLSPem != "" && s.cfg.TLSKey == "") || (s.cfg.TLSPem == "" && s.cfg.TLSKey != "") {
		return errors.New("(*Server).Listen() error: config.TLSPem or config.TLSKey doesn't have value")
	}

	switch s.cfg.TLSPem != "" && s.cfg.TLSKey != "" {
	case true:
		err = s.a.ListenAndServeTLS(s.cfg.TLSPem, s.cfg.TLSKey)

	case false:
		err = s.a.ListenAndServe()
	}

	return err
}

func (s *Server) GracefulShutdown(connectionsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigint

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.a.Shutdown(ctx); err != nil {
		s.log.Fatal(err.Error())
	}

	connectionsClosed <- struct{}{}
}
