package controllers

import "github.com/gin-gonic/gin"

// Handle404 handles 404 errors
func Handle404(c *gin.Context) {
	c.AbortWithStatusJSON(404, gin.H{
		"success": false,
		"error":   "Page not found",
	})
}
