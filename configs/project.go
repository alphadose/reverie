package configs

import "github.com/gin-gonic/gin"

var (
	// Project holds the main configuration for the entire project
	Project = getConfiguration()

	// MongoConfig is the configuration for MongoDB
	MongoConfig = Project.Mongo

	// AdminConfig is the configuration for default Gasper admin
	AdminConfig = Project.Admin
)

func init() {
	if Project.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
