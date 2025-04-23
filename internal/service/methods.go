package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"storage/internal/models"
	"strings"
	"time"
)

func WithDB(db updater) func(*Client) {
	return func(c *Client) {
		c.db = db
	}
}

func WithUserService(service userer) func(*Client) {
	return func(c *Client) {
		c.userService = service
	}
}
func WithLogger(logger *zap.Logger) func(*Client) {
	return func(c *Client) {
		c.logs = logger
	}
}

func New(opts ...func(*Client)) (*Client, error) {
	var cfg config
	e := env.Parse(&cfg)
	if e != nil {
		return nil, errors.Join(e, errors.New("error parsing config app.New"))
	}

	c := new(Client)
	for _, o := range opts {
		o(c)
	}

	c.cfg = cfg
	return c, nil
}

func (c *Client) Access(ctx context.Context, userGuid uuid.UUID, ip string) (models.TokenPair, error) {
	syncToken := uuid.New()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, models.TokenClaims{
		UserID:        userGuid,
		ClientIP:      ip,
		SyncTokenUUID: syncToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.cfg.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessTokenString, err := accessToken.SignedString([]byte(c.cfg.JWTSecret))
	if err != nil {
		return models.TokenPair{}, err
	}

	randomPart := uuid.New().String()
	refreshTokenContent := strings.Join([]string{userGuid.String(), syncToken.String(), randomPart}, ":")

	hashedPart, err := bcrypt.GenerateFromPassword([]byte(randomPart), c.cfg.BcryptCost)
	if err != nil {
		return models.TokenPair{}, err
	}

	err = c.db.Store(ctx, userGuid, string(hashedPart), time.Now().Add(c.cfg.RefreshTokenExp))
	if err != nil {
		return models.TokenPair{}, err
	}

	refreshTokenString := base64.StdEncoding.EncodeToString([]byte(refreshTokenContent))

	return models.TokenPair{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil

}

func (c *Client) Refresh(ctx context.Context, splited []string, header string, ip string) (models.TokenPair, error) {
	userId, err := uuid.Parse(splited[0])
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	syncToken, err := uuid.Parse(splited[1])
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	savedHash, expAt, err := c.db.Get(ctx, userId)
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	if time.Now().After(expAt) {
		c.logs.Error("token expired")
		return models.TokenPair{}, errors.New("token expired")
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(splited[2]))
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	accessClaims, err := c.validateAccessToken(tokenStr)
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	if accessClaims.SyncTokenUUID != syncToken {
		c.logs.Error("access and refresh token do not match")
		return models.TokenPair{}, errors.New("invalid refresh token")
	}

	if accessClaims.ClientIP != ip {
		err = c.userService.SendAlert(ctx, userId)
		if err != nil {
			c.logs.Error(err.Error())
			return models.TokenPair{}, err
		}

		c.logs.Error("invalid refresh token")
		return models.TokenPair{}, errors.New("invalid refresh token")
	}

	newTokens, err := c.Access(ctx, userId, ip)
	if err != nil {
		c.logs.Error(err.Error())
		return models.TokenPair{}, err
	}

	return newTokens, nil
}

func (c *Client) validateAccessToken(tokenString string) (*models.TokenClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(
		tokenString,
		&models.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(c.cfg.JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return nil, errors.New("invalid token structure")
	}

	return claims, nil
}
