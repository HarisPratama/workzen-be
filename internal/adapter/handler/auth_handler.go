package handler

import (
	"time"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	validatorLib "workzen-be/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResp response.ErrorResponseDefault
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func (a *authHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	resp := response.SuccessAuthResponse{}

	if refreshToken == "" {
		code = "[HANDLER] RefreshToken - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "unauthorized"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	newAccessToken, _, err := a.authService.RefreshToken(c.Context(), refreshToken)

	if err != nil {
		code = "[HANDLER] RefreshToken - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid refresh token"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	now := time.Now().Local()
	expiresAt := now.Add(15 * time.Minute)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		Domain:   ".workzen.web.id", // ← TAMBAH INI
		HTTPOnly: true,
		Secure:   true, // true jika HTTPS
		SameSite: "None",
		Expires:  expiresAt, // seconds
	})

	resp.Meta.Status = true
	resp.Meta.Message = "Refresh successful"

	return c.JSON(&resp)
}

func (a *authHandler) Logout(c *fiber.Ctx) error {
	resp := response.SuccessAuthResponse{}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   ".workzen.web.id", // ← TAMBAH INI
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   ".workzen.web.id", // ← TAMBAH INI
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
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
		errorResp.Meta.Message = err.Error()

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

		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	now := time.Now().Local()
	expiresAt := now.Add(15 * time.Minute)
	refreshExpires := now.Add(time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		Domain:   ".workzen.web.id", // ← TAMBAH INI
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  expiresAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		Domain:   ".workzen.web.id", // ← TAMBAH INI
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  refreshExpires,
	})

	resp.Meta.Status = true
	resp.Meta.Message = "Login successful"

	return c.JSON(&resp)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}
