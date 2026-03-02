package entity

import "github.com/golang-jwt/jwt/v5"

type JwtData struct {
	UserID float64 `json:"user_id"`
	Role   string  `json:"role"`
	jwt.RegisteredClaims
}
