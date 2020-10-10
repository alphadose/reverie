package configs

import "github.com/gin-gonic/gin"

var (
	// Project holds the main configuration for the entire project
	Project = getConfiguration()

	// MongoConfig is the configuration for MongoDB
	MongoConfig = Project.Mongo

	// AdminConfig is the configuration for default Reverie admin
	AdminConfig = Project.Admin

	// JWTConfig is the configuration for json web auth token
	JWTConfig = Project.JWT
)

func init() {
	if Project.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
