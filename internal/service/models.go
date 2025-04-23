package service

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	logs        *zap.Logger
	cfg         config
	db          updater
	userService userer
}

type updater interface {
	Get(ctx context.Context, userID uuid.UUID) (string, time.Time, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	Store(context.Context, uuid.UUID, string, time.Time) error
}

type userer interface {
	SendAlert(ctx context.Context, userID uuid.UUID) error
}

type config struct {
	AccessTokenExp  time.Duration `env:"ACCESS_TOKEN_EXP" envDefault:"1m"`
	RefreshTokenExp time.Duration `env:"REFRESH_TOKEN_EXP" envDefault:"1h"`
	BcryptCost      int           `env:"BCRYPT_COST" envDefault:"6"`
	JWTSecret       string        `env:"JWT_SECRET" envDefault:"secret"`
}
