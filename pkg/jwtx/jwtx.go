package jwtx

import "github.com/golang-jwt/jwt/v5"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type Claims struct {
	Role   string `json:"role"`
	UserID uint64 `json:"user_id,omitempty"`
	Name   string `json:"name,omitempty"`
	jwt.RegisteredClaims
}
