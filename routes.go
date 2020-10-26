package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	c "github.com/reverie/controllers"
	m "github.com/reverie/middlewares"
)

// PART LEFT: NOTIFICATIONS, Tracking remaining  requirements, Payment, Emails, Managing Vendor Inventory
// Request for increment/decrement of offer items to vendors

func newRouter() *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler: c.ErrorHandler,
		Prefork:      true,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Content-Length, Accept, Authorization, Cookie",
	}))

	auth := router.Group("/auth")
	{
		auth.Post("/login", c.Login)
		auth.Post("/register/client", c.RegisterClient)
		auth.Post("/register/vendor", c.RegisterVendor)
	}

	client := router.Group("/client", m.JWT, m.IsClient)
	{
		client.Get("", c.GetLoggedInUserInfo)
		client.Put("/password", c.UpdatePassword)
		client.Get("/post", c.FetchActivePostsByClient)
		client.Post("/post", c.CreatePost)
		client.Put("/post/:id", c.UpdatePost)
		// TODO: update vendor inventory after offer/  Debatable too constricting feature
		client.Patch("/post/:id/offer/:key/accept", c.AcceptOffer)
		// TODO: notify us when post is ongoing to handle end-to-end transactions such as logistics, payment etc
		client.Patch("/post/:id/activate", c.ActivatePost)
		client.Patch("/post/:id/deactivate", c.DeactivatePost)
		// TODO: update vendor inventory after completion/ Debatable too constricting feature
		// maybe reduce inventory on acceptance, seems right
		client.Patch("/post/:id/complete", c.MarkComplete)
	}

	vendor := router.Group("/vendor", m.JWT, m.IsVendor)
	{
		vendor.Get("", c.GetLoggedInUserInfo)
		vendor.Put("/inventory", c.UpdateInventory)
		vendor.Put("/password", c.UpdatePassword)
		vendor.Get("/post", c.FetchPostsByVendor)
		vendor.Put("/post/:id/offer", c.MakeOffer)
		vendor.Get("/post/offered", c.FetchOfferedPostsByVendor)
		vendor.Get("/post/contracted", c.FetchContractedPostsByVendor)
	}

	router.Use(c.Handle404)

	return router
}
