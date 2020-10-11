package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/reverie/controllers"
	m "github.com/reverie/middlewares"
)

func newRouter() http.Handler {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	router.NoRoute(c.Handle404)

	auth := router.Group("/auth")
	{
		auth.POST("/login", m.JWT.LoginHandler)
		auth.POST("/register/client", m.ValidateUserRegistration, c.RegisterClient)
		auth.POST("/register/vendor", m.ValidateUserRegistration, c.RegisterVendor)
		auth.GET("/refresh", m.JWT.RefreshHandler)
	}

	client := router.Group("/client")
	client.Use(m.JWT.MiddlewareFunc())
	client.Use(m.IsClient)
	{
		client.GET("", c.GetLoggedInUserInfo)
		client.PUT("/password", c.UpdatePassword)
		client.POST("/post", m.ValidatePostCreation, c.CreatePost)
	}

	vendor := router.Group("/vendor")
	vendor.Use(m.JWT.MiddlewareFunc())
	vendor.Use(m.IsVendor)
	{
		vendor.GET("", c.GetLoggedInUserInfo)
		vendor.PUT("/inventory", c.UpdateInventory)
		vendor.PUT("/password", c.UpdatePassword)
	}

	return router
}
