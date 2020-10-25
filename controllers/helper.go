package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/configs"
	"github.com/reverie/types"
)

// Handle404 handles 404 errors
func Handle404(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(types.M{
		types.Success: false,
		types.Error:   "Page not found",
	})
}

// ErrorHandler is the default error handler for the router
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	if configs.Project.Debug || code != fiber.StatusInternalServerError {
		return c.Status(code).JSON(types.M{
			types.Success: false,
			types.Error:   err.Error(),
		})
	}
	return c.Status(code).JSON(types.M{
		types.Success: false,
		types.Error:   "INTERNAL_SERVER_ERROR",
	})
}
