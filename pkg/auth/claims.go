package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"userid"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}
