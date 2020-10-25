package controllers

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// CreatePost creates a post requested by a client
func CreatePost(c *fiber.Ctx) error {
	post := &types.Post{}
	if err := c.BodyParser(post); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if result, err := validator.ValidateStruct(post); !result {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-1", utils.ErrFailedExtraction)
	}
	if err := post.Initialize(); err != nil {
		return utils.ServerError("Post-Controller-2", err)
	}
	post.SetOwner(claims.GetEmail())

	if _, err := mongo.CreatePost(post); err != nil {
		return utils.ServerError("Post-Controller-3", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}
