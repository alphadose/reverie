package controllers

import (
	"strconv"

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

// updatePostStatus updates the status of a post
func updatePostStatus(c *fiber.Ctx, status string) error {
	postID := c.Params("id")
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-8", utils.ErrFailedExtraction)
	}
	if err := mongo.UpdatePostStatus(postID, claims.GetEmail(), status); err != nil {
		return utils.ServerError("Post-Controller-9", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// ActivatePost intiates the post by marking its status as "ONGOING"
// No new offers can be made to this post
// This marks the start of the job defined in the post
func ActivatePost(c *fiber.Ctx) error {
	return updatePostStatus(c, types.ONGOING)
}

// DeactivatePost changes the post status from "ONGOING" to "OPEN"
// so that the client can accept new offers
func DeactivatePost(c *fiber.Ctx) error {
	return updatePostStatus(c, types.OPEN)
}

// MarkComplete marks the status of the post as "COMPLETED"
// Denotes the end of a job request
func MarkComplete(c *fiber.Ctx) error {
	return updatePostStatus(c, types.COMPLETED)
}

// UpdatePost updates the post by a client
// Can only update description, location and requirements
func UpdatePost(c *fiber.Ctx) error {
	postUpdate := &types.PostUpdate{}
	if err := c.BodyParser(postUpdate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !postUpdate.Location.Empty() {
		if result, err := validator.ValidateStruct(postUpdate.Location); !result {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := postUpdate.InitializeLocation(); err != nil {
			return utils.ServerError("Post-Controller-10", err)
		}
	}

	postID := c.Params("id")
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-11", utils.ErrFailedExtraction)
	}

	if err := mongo.UpdatePost(postID, claims.GetEmail(), postUpdate); err != nil {
		return utils.ServerError("Post-Controller-12", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// FetchOfferedPostsByVendor returns all open/ongoing posts the vendor has made an offer to
func FetchOfferedPostsByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-13", utils.ErrFailedExtraction)
	}
	offeredPosts, err := mongo.FetchOfferedPostsByVendor(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller14", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        offeredPosts,
	})
}

// FetchPostsByVendor returns all open posts
func FetchPostsByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-15", utils.ErrFailedExtraction)
	}
	page := c.Query("page", "0")
	pageNumber, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	openPosts, err := mongo.FetchPostsByVendor(claims.GetEmail(), pageNumber)
	if err != nil {
		return utils.ServerError("Post-Controller-16", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"page":        pageNumber,
		"data":        openPosts,
	})
}
