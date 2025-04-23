package http_transport

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"storage/internal/models"
)

type serviceInterface interface {
	Access(ctx context.Context, guid uuid.UUID, ip string) (models.TokenPair, error)
	Refresh(ctx context.Context, splited []string, header string, ip string) (models.TokenPair, error)
}

type Server struct {
	cfg config
	log *zap.Logger
	a   *http.Server

	service serviceInterface
}

type config struct {
	Host   string `env:"HOST" envDefault:":1235"`
	TLSKey string `env:"TLS_KEY" envDefault:""`
	TLSPem string `env:"TLS_PEM" envDefault:""`
}

type reqRefresh struct {
	Token string `json:"refresh_token"`
}
