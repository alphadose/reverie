package middlewares

import (
	"encoding/json"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/reverie/types"
)

// IsClient checks whether a user is a client or not
func IsClient(c *gin.Context) {
	user := ExtractClaims(c)
	if user.IsClient() {
		c.Next()
		return
	}
	c.JSON(403, gin.H{
		"success": false,
		"error":   "User is not a client",
	})
}

// IsVendor checks whether a user is a vendor or not
func IsVendor(c *gin.Context) {
	user := ExtractClaims(c)
	if user.IsVendor() {
		c.Next()
		return
	}
	c.JSON(403, gin.H{
		"success": false,
		"error":   "User is not a vendor",
	})
}

// ValidateUserRegistration validates the user registration request
func ValidateUserRegistration(c *gin.Context) {
	requestBody := getBodyFromContext(c)
	user := &types.User{}
	if err := json.Unmarshal(requestBody, user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result, err := validator.ValidateStruct(user); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.Next()
}
