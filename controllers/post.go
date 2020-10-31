package controllers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// CreatePost creates a post requested by a client
func CreatePost(c *fiber.Ctx) error {
	post := &types.Post{}
	if err := c.BodyParser(post); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if result, err := validator.ValidateStruct(post); !result {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-1", utils.ErrFailedExtraction)
	}
	if err := post.Initialize(); err != nil {
		return utils.ServerError("Post-Controller-2", err)
	}
	post.SetOwner(claims.GetEmail())

	if _, err := mongo.CreatePost(post); err != nil {
		return utils.ServerError("Post-Controller-3", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-4", utils.ErrFailedExtraction)
	}
	activePosts, err := mongo.FetchActivePostsByClient(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-5", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        activePosts,
	})
}

// MakeOffer adds/updates a vendor's offer to a post
func MakeOffer(c *fiber.Ctx) error {
	offer := &types.Inventory{}
	if err := c.BodyParser(offer); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-6", utils.ErrFailedExtraction)
	}

	postID := c.Params("id")
	requirements, err := mongo.FetchPostRequirements(postID)
	if err != nil {
		return utils.ServerError("Post-Controller", err)
	}

	vendorInventory, err := mongo.FetchVendorInventory(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller", err)
	}

	// Validate the offer made by the vendor
	offerValues := reflect.ValueOf(*offer)
	requirementValues := reflect.ValueOf(*requirements)
	vendorInventoryValues := reflect.ValueOf(*vendorInventory)

	for i := 0; i < offerValues.NumField(); i++ {
		offerVal := offerValues.Field(i).Int()
		reqVal := requirementValues.Field(i).Int()
		vendorInvVal := vendorInventoryValues.Field(i).Int()

		// Check if fields in the offer are negative or they exceed the post's requirements
		// or they exceed the vendor's inventory
		// If yes then return an error
		if offerVal < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values are negative")
		}
		if offerVal > reqVal {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the post's requirements")
		}
		if offerVal > vendorInvVal {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the vendor's current inventory limits")
		}
	}

	if err := mongo.UpdatePostOffers(postID, claims.GetEmail(), offer); err != nil {
		return utils.ServerError("Post-Controller-7", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// updatePostStatus updates the status of a post
func updatePostStatus(c *fiber.Ctx, status string) error {
	postID := c.Params("id")
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-8", utils.ErrFailedExtraction)
	}
	if err := mongo.UpdatePostStatus(postID, status); err != nil {
		return utils.ServerError("Post-Controller-9", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// ActivatePost intiates the post by marking its status as "ONGOING"
// No new offers can be made to this post
// This marks the start of the job defined in the post
func ActivatePost(c *fiber.Ctx) error {
	return updatePostStatus(c, types.ONGOING)
}

// DeactivatePost changes the post status from "ONGOING" to "OPEN"
// so that the client can accept new offers
func DeactivatePost(c *fiber.Ctx) error {
	return updatePostStatus(c, types.OPEN)
}

// MarkComplete marks the status of the post as "COMPLETED"
// Denotes the end of a job request
func MarkComplete(c *fiber.Ctx) error {
	// TODO : release all vendor inventories from accepted offers
	// postID := c.Params("id")

	// acceptedOffers, err := mongo.FetchPostAcceptedOffers(postID)
	return updatePostStatus(c, types.COMPLETED)
}

// UpdatePost updates the post by a client
// Can only update description, location and requirements
func UpdatePost(c *fiber.Ctx) error {
	postUpdate := &types.PostUpdate{}
	if err := c.BodyParser(postUpdate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !postUpdate.Location.Empty() {
		if result, err := validator.ValidateStruct(postUpdate.Location); !result {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := postUpdate.InitializeLocation(); err != nil {
			return utils.ServerError("Post-Controller-10", err)
		}
	}

	postID := c.Params("id")
	if err := mongo.UpdatePost(postID, postUpdate); err != nil {
		return utils.ServerError("Post-Controller-12", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// FetchOfferedPostsByVendor returns all open/ongoing posts the vendor has made an offer to
func FetchOfferedPostsByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-13", utils.ErrFailedExtraction)
	}
	offeredPosts, err := mongo.FetchOfferedPostsByVendor(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-14", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        offeredPosts,
	})
}

// inventoryCheckSet holds the correct inventory items
// used for detecting commodities which are not a part of our system ex:- space shuttles :3
var inventoryCheckSet *types.Set

// Initialize the set
func init() {
	inventoryCheckSet = types.NewSet()
	equipments := reflect.TypeOf(types.Inventory{})
	num := equipments.NumField()
	for i := 0; i < num; i++ {
		inventoryCheckSet.Add(equipments.Field(i).Name)
	}
}

// FetchPostsByVendor returns all open posts
func FetchPostsByVendor(c *fiber.Ctx) error {
	// Extract page number for pagination and validate
	page := c.Query("page", "0")
	pageNumber, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if pageNumber < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Page must be non-negative")
	}

	// Extract lookup items and validate
	lookupItems := strings.Split(c.Query("items"), ",")
	for _, item := range lookupItems {
		if !inventoryCheckSet.Contains(item) {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("%s is an invalid lookup item", item))
		}
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-13", utils.ErrFailedExtraction)
	}
	openPosts, err := mongo.FetchPostsByVendor(claims.GetEmail(), pageNumber, lookupItems)
	if err != nil {
		return utils.ServerError("Post-Controller-15", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"page":        pageNumber,
		"data":        openPosts,
	})
}

// FetchContractedPostsByVendor returns all posts in which the vendor's offer has been accepted
func FetchContractedPostsByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction)
	}
	contractedPosts, err := mongo.FetchContractedPostsByVendor(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-17", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        contractedPosts,
	})
}

// AcceptOffer accepts an offer made by a vendor on a post
// This operation is invoked by the client who is the owner of the post
// The param "offerKey" is key of the post holding the offer
// It is in the form of the vendor's email who made the offer with all "." replaced with "_"
// For Ex:- If the vendor's email is abc.2000@xyz.com the the key will be abc_2000@xyz_com
func AcceptOffer(c *fiber.Ctx) error {
	postID := c.Params("id")
	offerKey := c.Params("key")

	offers, acceptedOffers, requirements, err := mongo.FetchPostOffersAndRequirements(postID)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err)
	}

	// Check if offer exists
	offer, ok := offers[offerKey]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Offer key %s doesnt exist in post %s", offerKey, postID))
	}

	vendorEmail := strings.ReplaceAll(offerKey, "_", ".")
	vendorInventory, err := mongo.FetchVendorInventory(vendorEmail)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err)
	}

	// Check if offer exceeds post requirements or vendor's current inventory
	offerValues := reflect.ValueOf(offer)
	requirementValues := reflect.ValueOf(requirements)
	vendorInventoryValues := reflect.ValueOf(*vendorInventory)

	sanityChecker := make([]int64, offerValues.NumField())

	for i := 0; i < offerValues.NumField(); i++ {
		if offerValues.Field(i).Int() > vendorInventoryValues.Field(i).Int() {
			fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the vendor's current inventory limits")
		}
		sanityChecker[i] = requirementValues.Field(i).Int() - offerValues.Field(i).Int()
	}

	for _, acceptedOffer := range acceptedOffers {
		acceptedOfferValues := reflect.ValueOf(acceptedOffer)
		for i := 0; i < offerValues.NumField(); i++ {
			sanityChecker[i] -= acceptedOfferValues.Field(i).Int()
		}
	}

	for _, check := range sanityChecker {
		if check < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the post's requirements")
		}
	}

	if err := mongo.AcceptOffer(postID, offerKey, offer); err != nil {
		return utils.ServerError("Post-Controller-19", err)
	}

	if err := mongo.UpdateVendorInventoryOnAcceptance(vendorEmail, offer); err != nil {
		return utils.ServerError("Post-Controller-19", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}
