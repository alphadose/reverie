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

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-4", utils.ErrFailedExtraction)
	}
	activePosts, err := mongo.FetchActivePostsByClient(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-5", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        activePosts,
	})
}

// MakeOffer adds/updates a vendor's offer to a post
func MakeOffer(c *fiber.Ctx) error {
	postID := c.Params("id")
	offer := &types.Inventory{}
	if err := c.BodyParser(offer); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-6", utils.ErrFailedExtraction)
	}
	if err := mongo.UpdatePostOffers(postID, claims.GetEmail(), offer); err != nil {
		return utils.ServerError("Post-Controller-7", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}
