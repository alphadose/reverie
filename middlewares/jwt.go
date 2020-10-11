package middlewares

import (
	"errors"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/reverie/configs"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

var (
	errMissingCredentials   = errors.New("Missing Email or Password")
	errFailedAuthentication = errors.New("Incorrect Email or Password")
	// ErrFailedExtraction occurs when the request fails to extract claims from JWT
	ErrFailedExtraction = errors.New("Failed to extract JWT claims")
)

func authenticator(c *gin.Context) (interface{}, error) {
	auth := &types.Login{}
	if err := c.ShouldBind(auth); err != nil {
		return nil, errMissingCredentials
	}
	user, err := mongo.FetchSingleUser(auth.GetEmail())
	if err != nil || user == nil {
		return nil, errFailedAuthentication
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), auth.GetPassword()) {
		return nil, errFailedAuthentication
	}
	return user, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(*types.User); ok {
		return jwt.MapClaims{
			mongo.EmailKey:    user.Email,
			mongo.UsernameKey: user.Username,
			mongo.RoleKey:     user.Role,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	email, ok := claims[mongo.EmailKey].(string)
	if !ok {
		return nil
	}
	username, ok := claims[mongo.UsernameKey].(string)
	if !ok {
		return nil
	}
	role, ok := claims[mongo.RoleKey].(string)
	if !ok {
		return nil
	}
	return &types.User{
		Email:    email,
		Username: username,
		Role:     role,
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	_, ok := data.(*types.User)
	return ok
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}

// JWT handles the auth through JWT token
var JWT = &jwt.GinJWTMiddleware{
	Realm:           "Reverie",
	Key:             []byte(configs.JWTConfig.Secret),
	Timeout:         configs.JWTConfig.Timeout * time.Second,
	MaxRefresh:      configs.JWTConfig.MaxRefresh * time.Second,
	TokenLookup:     "header: Authorization",
	TokenHeadName:   "Bearer",
	TimeFunc:        time.Now,
	Authenticator:   authenticator,
	PayloadFunc:     payloadFunc,
	IdentityHandler: identityHandler,
	Authorizator:    authorizator,
	Unauthorized:    unauthorized,
}

// ExtractClaims takes the gin context and returns the User
// IMPORTANT: claims shall only provide the email, username and the role of a user
func ExtractClaims(c *gin.Context) *types.User {
	user, success := JWT.IdentityHandler(c).(*types.User)
	if !success {
		return nil
	}
	return user
}

func init() {
	// This keeps the middleware in check if the configuration is correct
	// Prevents runtime errors
	if err := JWT.MiddlewareInit(); err != nil {
		utils.Log("Master-JWT-1", "Failed to initialize JWT middleware", utils.ErrorTAG)
		utils.LogError("Master-JWT-2", err)
		os.Exit(1)
	}
}
