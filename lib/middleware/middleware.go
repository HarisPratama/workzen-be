package middleware

import (
	"fmt"
	"strings"
	"workzen-be/config"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/lib/auth"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
	CheckCookieToken() fiber.Handler
	RequireRole(roles ...string) fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

func (o Options) CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Authorization required"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)

		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid access token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func (o Options) CheckCookieToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault

		token := c.Cookies("access_token")

		if token == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Authorization required"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		claims, err := o.authJwt.VerifyAccessToken(token)

		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid access token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func (o Options) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault

		claims, ok := c.Locals("user").(*entity.JwtData)
		if !ok {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid token claims"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		role := claims.Role
		fmt.Println(role)

		if role == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Unauthorized access"
			return c.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		for _, allowed := range roles {
			if role == allowed {
				return c.Next()
			}
		}

		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Access forbidden"
		return c.Status(fiber.StatusForbidden).JSON(errorResponse)
	}
}

func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}
