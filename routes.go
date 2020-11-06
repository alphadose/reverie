package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/reverie/configs"
	c "github.com/reverie/controllers"
	m "github.com/reverie/middlewares"
)

// PART LEFT: NOTIFICATIONS, Payment, Emails
// Request for increment/decrement/removal of offer items to vendors (to be handled via notifications)
// Store JWT in local storage in frontend

// update post timestamp
// how to handle duplicate post creation ?
// Encrypt vendor emails in post offers and accepted offers ? (low priority)

func newRouter() *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler: c.ErrorHandler,
		Prefork:      !configs.Project.Debug,
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

		// Actions which only the owner of a post can perform
		postOwner := client.Group("/post/:id", m.IsPostOwner)
		{
			postOwner.Put("", c.UpdatePost)
			postOwner.Patch("/offer/:key/accept", c.AcceptOffer)
			postOwner.Delete("/offer/:key/reject", c.RejectOffer)

			// TODO: notify us when post is ongoing to handle end-to-end transactions such as logistics, payment etc
			postOwner.Patch("/activate", c.ActivatePost)
			postOwner.Patch("/deactivate", c.DeactivatePost)
			// TODO: make clients/vendors fill a survey after completion?
			postOwner.Patch("/complete", c.MarkComplete)
		}
	}

	vendor := router.Group("/vendor", m.JWT, m.IsVendor)
	{
		vendor.Get("", c.GetLoggedInUserInfo)
		vendor.Put("/inventory", c.InitializeInventory) // Restrict this, should only happen on our authorization
		vendor.Put("/password", c.UpdatePassword)
		vendor.Get("/post", c.FetchPostsByVendor)
		// TODO: notify us so that we can contact the client directly in case he doesnt use the app
		// Always make sure to update the entire body i.e the new body will be the new offer entirely (it replaces the old body, not updates it)
		vendor.Put("/post/:id/offer", c.MakeOffer)
		vendor.Delete("/post/:id/retract", c.RetractOffer)

		vendor.Get("/post/offered", c.FetchOfferedPostsByVendor)
		vendor.Get("/post/contracted", c.FetchContractedPostsByVendor)
	}

	router.Use(c.Handle404)

	return router
}
