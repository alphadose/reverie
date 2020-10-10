package middlewares

import (
	"encoding/json"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/reverie/types"
)

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
