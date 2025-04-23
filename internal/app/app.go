package app

import (
	"errors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"storage/internal/db"
	"storage/internal/service"
	"storage/internal/transport/http_transport"
	"storage/internal/userService"
)

func Start() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	dbClient, err := db.New()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbClient.Close()

	userSer, err := userService.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	s, err := service.New(service.WithDB(dbClient), service.WithUserService(userSer), service.WithLogger(logger))
	if err != nil {
		logger.Fatal(err.Error())
	}

	transport, e := http_transport.New(s, logger)
	if e != nil {
		log.Fatal(e)
	}

	connectionsClosed := make(chan struct{})
	go transport.GracefulShutdown(connectionsClosed)

	logger.Info("Application started")

	if err := transport.Listen(); errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err.Error())
	} else {
		<-connectionsClosed
	}
}
