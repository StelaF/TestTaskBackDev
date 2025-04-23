package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	UserID        uuid.UUID `json:"user_id"`
	ClientIP      string    `json:"client_ip"`
	SyncTokenUUID uuid.UUID `json:"token_uuid"`
	jwt.RegisteredClaims
}
