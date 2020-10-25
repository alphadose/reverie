package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/reverie/configs"
	"github.com/reverie/types"
)

// The main error handler for JWT authentication
func authErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(types.M{
		types.Success: false,
		types.Error:   err.Error(),
	})
}

// JWT handles the auth through JWT token
var JWT = jwtware.New(jwtware.Config{
	SigningKey:   []byte(configs.JWTConfig.Secret),
	ErrorHandler: authErrorHandler,
})
