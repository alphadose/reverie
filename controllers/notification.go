package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// FetchNotifications returns all notifications for a user
func FetchNotifications(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction)
	}
	// Extract page number for pagination and validate
	page := c.Query("page", "0")
	pageNumber, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	notifications, err := mongo.FetchNotifications(claims.GetEmail(), pageNumber)
	if err != nil {
		return utils.ServerError("Post-Controller-19", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"page":        pageNumber,
		"data":        notifications,
	})
}

// ReadNotification marks a notification as "Read"
func ReadNotification(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction)
	}

	notificationID := c.Params("id")
	if err := mongo.MarkRead(notificationID, claims.GetEmail()); err != nil {
		return utils.ServerError("Post-Controller-19", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}
