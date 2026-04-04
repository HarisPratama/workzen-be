package service

import (
	"context"
	"errors"
	"strconv"
	"time"
	"workzen-be/config"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/lib/auth"
	"workzen-be/lib/conv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var err error
var code string

var (
	ErrUserNotFound    = errors.New("email not registered")
	ErrInvalidPassword = errors.New("incorrect password")
	ErrAccountInactive = errors.New("account is inactive, please contact administrator")
	ErrTokenGeneration = errors.New("failed to generate authentication token")
	ErrInvalidRefresh  = errors.New("refresh token is invalid or expired")
)

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
		return "", 0, ErrInvalidRefresh
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
	if err != nil {
		code = "[SERVICE] RefreshToken - 2"
		log.Errorw(code, err)
		return "", 0, ErrTokenGeneration
	}

	return accessToken, expiresAt, nil
}

func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if result.Status != "ACTIVE" {
		code = "[SERVICE] GetUserByEmail - 2"
		log.Errorw(code, "account inactive", "email", req.Email)
		return nil, ErrAccountInactive
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, "Invalid Password")
		return nil, ErrInvalidPassword
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
			ID:        strconv.FormatInt(result.ID, 10),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 4"
		log.Errorw(code, err)
		return nil, ErrTokenGeneration
	}

	refreshToken, _, err := a.jwtToken.GenerateRefreshToken(&jwtData)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 5"
		log.Errorw(code, err)
		return nil, ErrTokenGeneration
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
