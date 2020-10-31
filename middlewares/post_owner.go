package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/models/mongo"
	"github.com/reverie/utils"
)

// IsPostOwner checks whether the client is the owner of the post or not
func IsPostOwner(c *fiber.Ctx) error {
	postID := c.Params("id")
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Middleware-Validator-1", utils.ErrFailedExtraction)
	}

	owner, err := mongo.IsPostOwner(postID, claims.GetEmail())
	if err != nil {
		return utils.ServerError("Middleware-Validator-1", err)
	}

	if !owner {
		return fiber.NewError(fiber.StatusForbidden, "Client is not the owner of the post")
	}
	return c.Next()
}
