package controllers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

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
		return utils.ServerError("Post-Controller-1", utils.ErrFailedExtraction, c)
	}
	if err := post.Initialize(); err != nil {
		return utils.ServerError("Post-Controller-2", err, c)
	}
	post.SetOwner(claims.GetEmail())
	post.SetOwnerName(claims.GetName())
	id, err := mongo.CreatePost(post)
	if err != nil {
		return utils.ServerError("Post-Controller-3", err, c)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"_id":         id,
	})
}

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-4", utils.ErrFailedExtraction, c)
	}
	activePosts, err := mongo.FetchActivePostsByClient(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-5", err, c)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        activePosts,
	})
}

// FetchSinglePostByClient returns a single post given its id
func FetchSinglePostByClient(c *fiber.Ctx) error {
	postID := c.Params("id")

	post, err := mongo.FetchSinglePostByClient(postID)
	if err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

// MakeOffer adds/updates a vendor's offer to a post
func MakeOffer(c *fiber.Ctx) error {
	rate, err := strconv.ParseFloat(c.Params("rate"), 64)
	if err != nil {
		return utils.ServerError("Post-Controller-6", err, c)
	}

	offer := &types.Inventory{}
	if err := c.BodyParser(offer); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-6", utils.ErrFailedExtraction, c)
	}

	postID := utils.ImmutableString(c.Params("id"))
	status, requirements, err := mongo.FetchPostRequirementsAndStatus(postID)

	if status != types.OPEN {
		return fiber.NewError(fiber.StatusForbidden, "Offers can be made only to OPEN posts")
	}

	if err != nil {
		return utils.ServerError("Post-Controller", err, c)
	}

	vendorInventory, err := mongo.FetchVendorInventory(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller", err, c)
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

	if err := mongo.UpdatePostOffers(postID, claims.GetEmail(), types.Offer{
		Name:    claims.GetName(),
		Created: time.Now().Unix(),
		Content: *offer,
		Rate:    rate,
	}); err != nil {
		return utils.ServerError("Post-Controller-7", err, c)
	}

	go mongo.NotifyClient(postID, claims.GetName()+" made an offer to your post %s")

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// RetractOffer removes a vendor's offer from a post
func RetractOffer(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-6", utils.ErrFailedExtraction, c)
	}

	postID := utils.ImmutableString(c.Params("id"))

	if err := mongo.RetractPostOffer(postID, claims.GetEmail()); err != nil {
		return utils.ServerError("Post-Controller-7", err, c)
	}

	go mongo.NotifyClient(postID, claims.GetName()+" retracted his offer from your post %s")

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// updatePostStatus updates the status of a post
func updatePostStatus(c *fiber.Ctx, status string) error {
	postID := utils.ImmutableString(c.Params("id"))
	if err := mongo.UpdatePostStatus(postID, status); err != nil {
		return utils.ServerError("Post-Controller-9", err, c)
	}

	// Notify all vendors whose offers have been accepted
	go mongo.BulkNotifyVendors(postID, status)

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

// DeletePost changes the post status to "DELETED"
func DeletePost(c *fiber.Ctx) error {
	return updatePostStatus(c, types.DELETED)
}

// MarkComplete marks the status of the post as "COMPLETED"
// Denotes the end of a job request
func MarkComplete(c *fiber.Ctx) error {
	postID := c.Params("id")
	acceptedOffers, status, err := mongo.FetchPostAcceptedOffersAndStatus(postID)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	if status != types.ONGOING {
		return fiber.NewError(fiber.StatusForbidden, "Only ONGOING posts can be marked completed")
	}

	if err := mongo.ReleaseVendorInventories(acceptedOffers); err != nil {
		return utils.ServerError("kekw", err, c)
	}
	return updatePostStatus(c, types.COMPLETED)
}

// UpdatePost updates the post by a client
// Can only update description, location and requirements
func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	status, err := mongo.FetchPostStatus(postID)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	if status != types.OPEN {
		return fiber.NewError(fiber.StatusForbidden, "Only OPEN posts can be updated")
	}

	postUpdate := &types.PostUpdate{}
	if err := c.BodyParser(postUpdate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if postUpdate.Location != nil {
		if result, err := validator.ValidateStruct(postUpdate.Location); !result {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := postUpdate.InitializeLocation(); err != nil {
			return utils.ServerError("Post-Controller-10", err, c)
		}
	}

	if err := mongo.UpdatePost(postID, postUpdate); err != nil {
		return utils.ServerError("Post-Controller-12", err, c)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// FetchOfferedPostsByVendor returns all open/ongoing posts the vendor has made an offer to
func FetchOfferedPostsByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-13", utils.ErrFailedExtraction, c)
	}
	offeredPosts, err := mongo.FetchOfferedPostsByVendor(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-14", err, c)
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
		return utils.ServerError("Post-Controller-13", utils.ErrFailedExtraction, c)
	}
	openPosts, err := mongo.FetchPostsByVendor(claims.GetEmail(), pageNumber, lookupItems)
	if err != nil {
		return utils.ServerError("Post-Controller-15", err, c)
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
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction, c)
	}
	contractedPosts, err := mongo.FetchContractedPostsByVendor(claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-17", err, c)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"data":        contractedPosts,
	})
}

// FetchSinglePostByVendor returns a single post given its id
func FetchSinglePostByVendor(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("Post-Controller-16", utils.ErrFailedExtraction, c)
	}
	postID := c.Params("id")

	post, err := mongo.FetchSinglePostByVendor(postID, claims.GetEmail())
	if err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

// AcceptOffer accepts an offer made by a vendor on a post
// This operation is invoked by the client who is the owner of the post
// The param "offerKey" is key holding the offer in the post
// It is the vendor's email address encrypted with AES-256
func AcceptOffer(c *fiber.Ctx) error {
	postID := utils.ImmutableString(c.Params("id"))
	offerKey := c.Params("key")

	status, offers, requirements, err := mongo.FetchPostOffersAndRequirementsAndStatus(postID)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err, c)
	}

	if status != types.OPEN {
		return fiber.NewError(fiber.StatusForbidden, "Offers can be accepted only on OPEN posts")
	}

	// Check if offer exists
	offer, ok := offers[offerKey]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Offer key %s doesnt exist in post %s", offerKey, postID))
	}

	vendorEmail, err := utils.Decrypt(offerKey)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	vendorInventory, err := mongo.FetchVendorInventory(vendorEmail)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err, c)
	}

	// Check if offer exceeds post requirements or vendor's current inventory
	offerValues := reflect.ValueOf(offer.Content)
	requirementValues := reflect.ValueOf(requirements)
	vendorInventoryValues := reflect.ValueOf(*vendorInventory)

	sanityChecker := make([]int64, offerValues.NumField())

	for i := 0; i < offerValues.NumField(); i++ {
		if offerValues.Field(i).Int() > vendorInventoryValues.Field(i).Int() {
			fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the vendor's current inventory limits")
		}
		sanityChecker[i] = requirementValues.Field(i).Int() - offerValues.Field(i).Int()
	}

	for _, check := range sanityChecker {
		if check < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the post's requirements")
		}
	}

	if err := mongo.AcceptOffer(postID, offerKey, offer); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	if err := mongo.UpdateVendorInventoryOnAcceptance(vendorEmail, offer.Content); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	// Notify vendor
	go mongo.NotifyVendorOnAcceptance(postID, vendorEmail)

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// RejectAcceptedOffer removes an accepted offer by a client and adds the offer's contents back to the vendor's inventory
func RejectAcceptedOffer(c *fiber.Ctx) error {
	postID := utils.ImmutableString(c.Params("id"))
	offerKey := c.Params("key")

	acceptedOffers, status, err := mongo.FetchPostAcceptedOffersAndStatus(postID)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	if status != types.OPEN {
		return fiber.NewError(fiber.StatusForbidden, "Accepted Offers can be rejected only on OPEN posts")
	}

	// Check if offer exists
	offer, ok := acceptedOffers[offerKey]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Accepted Offer key %s doesnt exist in post %s", offerKey, postID))
	}

	vendorEmail, err := utils.Decrypt(offerKey)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	if err := mongo.RejectAcceptedOffer(postID, offerKey, offer.Content); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	if err := mongo.ReleaseSingleVendorInventory(vendorEmail, offer.Content); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	// Notify vendor
	go mongo.NotifyVendorOnRejection(postID, vendorEmail)

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// RejectPendingOffer removes a pending offer by a client
func RejectPendingOffer(c *fiber.Ctx) error {
	postID := utils.ImmutableString(c.Params("id"))
	offerKey := c.Params("key")

	vendorEmail, err := utils.Decrypt(offerKey)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	if err := mongo.RejectPendingOffer(postID, offerKey); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	// Notify vendor
	go mongo.NotifyVendorOnRejection(postID, vendorEmail)

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// RequestOfferChange is a request for change by the client on an offer on his post
// This operation is invoked by the client who is the owner of the post
// Ex:- Suppose a vendor has offered 3 trucks but the client only wants 2 trucks
// In this case the client requests the vendor to change his offer to 2 trucks
// The param "offerKey" is key holding the offer in the post
// The body "offerChange" holds the final offer desired by the client
// It is the vendor's email address encrypted with AES-256
func RequestOfferChange(c *fiber.Ctx) error {
	postID := c.Params("id")
	offerKey := c.Params("key")

	offerChange := &types.Inventory{}
	if err := c.BodyParser(offerChange); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, offers, requirements, err := mongo.FetchPostOffersAndRequirementsAndStatus(postID)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err, c)
	}

	if status != types.OPEN {
		return fiber.NewError(fiber.StatusForbidden, "Offers can be accepted only on OPEN posts")
	}

	// Check if offer exists
	_, ok := offers[offerKey]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Offer key %s doesnt exist in post %s", offerKey, postID))
	}

	vendorEmail, err := utils.Decrypt(offerKey)
	if err != nil {
		return utils.ServerError("kekw", err, c)
	}

	vendorInventory, err := mongo.FetchVendorInventory(vendorEmail)
	if err != nil {
		return utils.ServerError("Post-Controller-18", err, c)
	}

	// Check if offer exceeds post requirements or vendor's current inventory
	offerValues := reflect.ValueOf(offerChange)
	requirementValues := reflect.ValueOf(requirements)
	vendorInventoryValues := reflect.ValueOf(*vendorInventory)

	sanityChecker := make([]int64, offerValues.NumField())

	for i := 0; i < offerValues.NumField(); i++ {
		if offerValues.Field(i).Int() > vendorInventoryValues.Field(i).Int() {
			fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the vendor's current inventory limits")
		}
		sanityChecker[i] = requirementValues.Field(i).Int() - offerValues.Field(i).Int()
	}

	for _, check := range sanityChecker {
		if check < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Offer values exceed the post's requirements")
		}
	}

	if err := mongo.NotifyOfferChangeToVendor(postID, vendorEmail, offerChange); err != nil {
		return utils.ServerError("Post-Controller-19", err, c)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}
