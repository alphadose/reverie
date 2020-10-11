package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/reverie/middlewares"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// CreatePost creates a post requested by a client
func CreatePost(c *gin.Context) {
	post := &types.Post{}
	if err := c.BindJSON(post); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse("Post-Controller-1", c, middlewares.ErrFailedExtraction)
		return
	}
	if err := post.Initialize(); err != nil {
		utils.SendServerErrorResponse("Post-Controller-2", c, err)
		return
	}
	post.SetOwner(claims.GetEmail())

	if _, err := mongo.CreatePost(post); err != nil {
		utils.SendServerErrorResponse("Post-Controller-3", c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "post created",
	})
}
