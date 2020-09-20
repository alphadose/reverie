package configs

import "github.com/gin-gonic/gin"

// Project holds the main configuration for the entire project
var Project = getConfiguration()

func init() {
	if Project.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
