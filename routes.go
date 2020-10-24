package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func newRouter() *fiber.App {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Content-Length, Accept, Authorization, Cookie",
	}))

	// auth := router.Group("/auth")
	// {
	// 	auth.Post("/login", m.JWT.LoginHandler)
	// 	auth.Post("/register/client", m.ValidateUserRegistration, c.RegisterClient)
	// 	auth.Post("/register/vendor", m.ValidateUserRegistration, c.RegisterVendor)
	// 	auth.Get("/refresh", m.JWT.RefreshHandler)
	// }

	// client := router.Group("/client")
	// client.Use(m.JWT.MiddlewareFunc())
	// client.Use(m.IsClient)
	// {
	// 	client.Get("", c.GetLoggedInUserInfo)
	// 	client.Put("/password", c.UpdatePassword)
	// 	client.Post("/post", m.ValidatePostCreation, c.CreatePost)
	// }

	// vendor := router.Group("/vendor")
	// vendor.Use(m.JWT.MiddlewareFunc())
	// vendor.Use(m.IsVendor)
	// {
	// 	vendor.Get("", c.GetLoggedInUserInfo)
	// 	vendor.Put("/inventory", c.UpdateInventory)
	// 	vendor.Put("/password", c.UpdatePassword)
	// }

	// router.Use(c.Handle404)

	return router
}
