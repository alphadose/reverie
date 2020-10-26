package mongo

import (
	"fmt"
	"strings"

	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// postCollectionKey is the collection for all job posts
	postCollectionKey = "posts"

	// postRequirementsKey is the key denoting the requirements for a post
	postRequirementsKey = "requirements"

	// postOffersKey is the key denoting the offers made to a post by a vendor
	postOffersKey = "offers"

	// postLocationKey is the key denoting the location of a job request
	postLocationKey = "location"

	// postOwnerKey is the key holding the owner email of a post
	postOwnerKey = "owner"

	// postStatusKey is the key holding the status of a post
	postStatusKey = "status"

	// createdKey is the key denoting the timestamp of creation of a job request
	createdKey = "created"

	// updatedKey is the key denoting the timestamp at which the job request was last updated
	updatedKey = "updated"
)

var postCollection = db.Collection(postCollectionKey)

// convert "." to "_" for storing in mongoDB
func processEmail(email string) string {
	return strings.ReplaceAll(email, ".", "_")
}

// CreatePost is an abstraction over InsertOne which inserts a post
func CreatePost(post *types.Post) (interface{}, error) {
	return insertOne(postCollection, post)
}

// UpdatePost updates a post by a client
func UpdatePost(postID, clientEmail string, post *types.PostUpdate) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey:   docID,
		postOwnerKey: clientEmail,
	}
	return updateOne(postCollection, filter, post)
}

// UpdatePostOffers adds/updates an offer to a post
func UpdatePostOffers(postID, vendorEmail string, offer *types.Inventory) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	updatePayload := types.M{
		fmt.Sprintf("%s.%s", postOffersKey, processEmail(vendorEmail)): offer,
	}
	return updateOne(postCollection, filter, updatePayload, options.FindOneAndUpdate().SetUpsert(true))
}

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(email string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postOwnerKey: email,
		"$or": []types.M{
			{postStatusKey: types.OPEN},
			{postStatusKey: types.ONGOING},
		},
	})
}

// UpdatePostStatus updates the status of the post
func UpdatePostStatus(postID, clientEmail, newStatus string) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey:   docID,
		postOwnerKey: clientEmail,
	}
	updatePayload := types.M{
		postStatusKey: newStatus,
	}
	return updateOne(postCollection, filter, updatePayload)
}

// FetchOfferedPostsByVendor returns all open/ongoing posts the vendor has made an offer to
func FetchOfferedPostsByVendor(vendorEmail string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.OPEN,
		fmt.Sprintf("%s.%s", postOffersKey, processEmail(vendorEmail)): types.M{
			"$exists": true,
		},
	})
}
