package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/types"
)

// ServerError sends internal server error messages
func ServerError(logContext string, err error) error {
	LogError(logContext, err)
	return err
}

var (
	// ErrFailedExtraction occurs when the request fails to extract claims from JWT
	ErrFailedExtraction = errors.New("Failed to extract JWT claims")
)

// ExtractClaims takes the fiber context and returns the User
// IMPORTANT: claims shall only provide the email, username and the role of a user
func ExtractClaims(c *fiber.Ctx) *types.Claims {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email, ok := claims[types.EmailKey].(string)
	if !ok {
		return nil
	}
	username, ok := claims[types.UsernameKey].(string)
	if !ok {
		return nil
	}
	role, ok := claims[types.RoleKey].(string)
	if !ok {
		return nil
	}
	return &types.Claims{
		Email:    email,
		Username: username,
		Role:     role,
	}
}
