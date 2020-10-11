package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/reverie/middlewares"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// registerUser handles registration of new users
func registerUser(c *gin.Context, role string) {
	user := &types.User{}
	if err := c.BindJSON(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	filter := types.M{mongo.EmailKey: user.Email}
	count, err := mongo.CountUsers(filter)
	if err != nil {
		utils.LogError("User-Controller-1", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "email already registered",
		})
		return
	}

	hashedPass, err := utils.HashPassword(user.GetPassword())
	if err != nil {
		utils.LogError("User-Controller-2", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetPassword(hashedPass)
	user.SetRole(role)

	if _, err = mongo.RegisterUser(user); err != nil {
		utils.LogError("User-Controller-3", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "user created",
	})
}

// RegisterClient registers a new client
func RegisterClient(c *gin.Context) {
	registerUser(c, types.Client)
}

// RegisterVendor registers a new vendor
func RegisterVendor(c *gin.Context) {
	registerUser(c, types.Vendor)
}

// GetUserInfo gets info regarding particular user
func GetUserInfo(c *gin.Context) {
	user, err := mongo.FetchSingleUserWithoutPassword(c.Param("user"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "No such user exists",
			})
			return
		}
		utils.LogError("User-Controller-4", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetSuccess(true)
	c.JSON(200, user)
}

// GetLoggedInUserInfo returns info regarding the current logged in user
func GetLoggedInUserInfo(c *gin.Context) {
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	user, err := mongo.FetchSingleUserWithoutPassword(claims.GetEmail())
	if err != nil {
		utils.LogError("User-Controller-5", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetSuccess(true)
	c.JSON(200, user)
}

// UpdatePassword updates the password of a user
func UpdatePassword(c *gin.Context) {
	passwordUpdate := &types.PasswordUpdate{}
	if err := c.ShouldBind(passwordUpdate); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.LogError("User-Controller-6", middlewares.ErrFailedExtraction, c)
		utils.SendServerErrorResponse(c, middlewares.ErrFailedExtraction)
		return
	}
	user, err := mongo.FetchSingleUser(claims.GetEmail())
	if err != nil {
		utils.LogError("User-Controller-7", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), passwordUpdate.GetOldPassword()) {
		c.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   "old password is invalid",
		})
		return
	}
	hashedPass, err := utils.HashPassword(passwordUpdate.GetNewPassword())
	if err != nil {
		utils.LogError("User-Controller-8", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	err = mongo.UpdateUser(
		types.M{mongo.EmailKey: user.GetEmail()},
		types.M{mongo.PasswordKey: hashedPass},
	)
	if err != nil {
		utils.LogError("User-Controller-9", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "password updated",
	})
}

// DeleteUser deletes the user from database
func DeleteUser(c *gin.Context) {
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.LogError("User-Controller-10", middlewares.ErrFailedExtraction, c)
		utils.SendServerErrorResponse(c, middlewares.ErrFailedExtraction)
		return
	}
	filter := types.M{
		mongo.EmailKey: claims.GetEmail(),
	}
	updatePayload := types.M{
		"deleted": true,
	}
	err := mongo.UpdateUser(filter, updatePayload)
	if err != nil {
		utils.LogError("User-Controller-11", err, c)
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "user deleted",
	})
}
