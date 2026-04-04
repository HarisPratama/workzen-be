package handler

import (
	"errors"
	"time"
	"workzen-be/config"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResp response.ErrorResponseDefault
var jakartaTZ, _ = time.LoadLocation("Asia/Jakarta")

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

func (a *authHandler) cookieConfig() (string, bool, string) {
	domain := a.cfg.App.CookieDomain
	secure := a.cfg.App.CookieSecure
	sameSite := a.cfg.App.CookieSameSite
	if sameSite == "" {
		sameSite = "Lax"
	}
	return domain, secure, sameSite
}

func (a *authHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	resp := response.SuccessAuthResponse{}

	if refreshToken == "" {
		code = "[HANDLER] RefreshToken - 1"
		log.Errorw(code, "refresh token cookie is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Session expired, please login again"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	newAccessToken, _, err := a.authService.RefreshToken(c.Context(), refreshToken)

	if err != nil {
		code = "[HANDLER] RefreshToken - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false

		if errors.Is(err, service.ErrInvalidRefresh) {
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}

		errorResp.Meta.Message = "Failed to refresh session"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	domain, secure, sameSite := a.cookieConfig()
	now := time.Now().In(jakartaTZ)
	expiresAt := now.Add(15 * time.Minute)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		Domain:   domain,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Expires:  expiresAt,
	})

	resp.Meta.Status = true
	resp.Meta.Message = "Refresh successful"

	return c.JSON(&resp)
}

func (a *authHandler) Logout(c *fiber.Ctx) error {
	resp := response.SuccessAuthResponse{}
	domain, secure, sameSite := a.cookieConfig()

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})

	resp.Meta.Status = true
	resp.Meta.Message = "Logout successful"

	return c.JSON(&resp)
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	resp := response.SuccessAuthResponse{}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body, please check your input format"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] Login - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(errorResp)
		case errors.Is(err, service.ErrInvalidPassword):
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		case errors.Is(err, service.ErrAccountInactive):
			return c.Status(fiber.StatusForbidden).JSON(errorResp)
		case errors.Is(err, service.ErrTokenGeneration):
			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		default:
			errorResp.Meta.Message = "An unexpected error occurred, please try again later"
			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	domain, secure, sameSite := a.cookieConfig()
	now := time.Now().In(jakartaTZ)
	expiresAt := now.Add(15 * time.Minute)
	refreshExpires := now.Add(time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		Domain:   domain,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Expires:  expiresAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		Domain:   domain,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Expires:  refreshExpires,
	})

	resp.Meta.Status = true
	resp.Meta.Message = "Login successful"

	return c.JSON(&resp)
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) AuthHandler {
	return &authHandler{
		authService: authService,
		cfg:         cfg,
	}
}
