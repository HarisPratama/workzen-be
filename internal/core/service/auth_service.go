package service

import (
	"bwanews/config"
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"bwanews/lib/auth"
	"bwanews/lib/conv"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var err error
var code string

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, int64, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       auth.Jwt
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (string, int64, error) {
	claims, err := a.jwtToken.VerifyAccessToken(refreshToken)

	if err != nil {
		code = "[SERVICE] RefreshToken - 1"
		log.Errorw(code, err)
		return "", 0, err
	}

	jwtData := entity.JwtData{
		UserID:   claims.UserID,
		Role:     claims.Role,
		TenantID: claims.TenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
			ID:        string(claims.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)

	return accessToken, expiresAt, err
}

func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[SERVICE] GetUserByEmail - 2"
		err = errors.New("invalid password")
		log.Errorw(code, "Invalid Password")
		return nil, err
	}

	var tenantID float64
	if result.TenantID != nil {
		tenantID = float64(*result.TenantID)
	}

	jwtData := entity.JwtData{
		UserID:   float64(result.ID),
		Role:     result.Role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
			ID:        string(result.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	refreshToken, _, err := a.jwtToken.GenerateRefreshToken(&jwtData)

	if err != nil {
		code = "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.AccessToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	return &resp, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{authRepository: authRepository, cfg: cfg, jwtToken: jwtToken}
}
