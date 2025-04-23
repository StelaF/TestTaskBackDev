package userService

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func New() (*User, error) {
	var cfg config
	e := env.Parse(&cfg)
	if e != nil {
		return nil, e
	}

	return &User{
		cfg: cfg,
	}, nil
}

func (u *User) SendAlert(ctx context.Context, userID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	logger.Info("Sending Alert On example@example.com uuid:" + userID.String())
	return nil
}
