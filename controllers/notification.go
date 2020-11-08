package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// FetchUnreadNotifications returns all unread notifications for a user
func FetchUnreadNotifications(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction)
	}

	notifications, err := mongo.FetchUnreadNotifications(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-19", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
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
