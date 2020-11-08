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
// Validate email via link during registration

// Need to make a rulebook for the support team
// Contents :-
// 1. For new accepted offer delete the pending one, then ask the vendor to make a new offer and accept the new one
// 2. For adding contents to existing accepted offer, ask the vendor to make an offer with the remainder, then accept this offer, the new offer gets merged with the existing accepted offer
// 3. No auditing, hence nothing can be done if accepted/pending offer is deleted (auditing required) or maybe just mark deleted ?

// TODO : refactor mongo code
// update post timestamp

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
		// TODO : add email verification here
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
			postOwner.Delete("", c.DeletePost)
			postOwner.Patch("/offer/:key/accept", c.AcceptOffer)
			postOwner.Delete("/offer/:key/reject-accepted", c.RejectAcceptedOffer)
			postOwner.Delete("/offer/:key/reject-pending", c.RejectPendingOffer)

			// TODO: notify us when post is ongoing to handle end-to-end transactions such as logistics, payment etc
			postOwner.Patch("/activate", c.ActivatePost)
			postOwner.Patch("/deactivate", c.DeactivatePost)

			// TODO: make clients/vendors fill a survey after completion?
			// Restrict this route? client hits this, then we get a mail and approve and then only the process gets completed
			// We shall hit the admin route
			// This will generate the payment invoice and mail the client
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
