package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/utils"
)

// IsClient checks whether a user is a client or not
func IsClient(c *fiber.Ctx) error {
	user := utils.ExtractClaims(c)
	if user == nil {
		return utils.ServerError("Middleware-Validator-1", utils.ErrFailedExtraction, c)
	}
	if user.IsClient() || user.IsAdmin() {
		return c.Next()
	}
	return fiber.NewError(fiber.StatusForbidden, "User is not a client")
}

// IsVendor checks whether a user is a vendor or not
func IsVendor(c *fiber.Ctx) error {
	user := utils.ExtractClaims(c)
	if user == nil {
		return utils.ServerError("Middleware-Validator-2", utils.ErrFailedExtraction, c)
	}
	if user.IsVendor() || user.IsAdmin() {
		return c.Next()
	}
	return fiber.NewError(fiber.StatusForbidden, "User is not a vendor")
}
