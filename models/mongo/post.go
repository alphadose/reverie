package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	// postAcceptedOffersKey is the key denoting the offers accepted on a post by a client
	postAcceptedOffersKey = "accepted_offers"

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

// concatenates 2 strings with "." in betweem
func concat(A, B string) string {
	return fmt.Sprintf("%s.%s", A, B)
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
		concat(postOffersKey, processEmail(vendorEmail)): offer,
	}
	return updateOne(postCollection, filter, updatePayload)
}

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(email string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postOwnerKey: email,
		postStatusKey: types.M{
			"$in": []string{types.OPEN, types.ONGOING},
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

// FetchPostsByVendor returns all open posts based on the vendor's inventory
// TODO: be sure to add to projections on addition of sensitive fields to posts
func FetchPostsByVendor(vendorEmail string, pageNumber int64, lookupItems []string) ([]types.M, error) {
	searchArray := make([]types.M, 0)
	for _, item := range lookupItems {
		searchArray = append(searchArray, types.M{
			concat(postRequirementsKey, item): types.M{
				"$gt": 0,
			},
		})
	}
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.OPEN,
		"$or":         searchArray,
		concat(postOffersKey, processEmail(vendorEmail)): types.M{
			"$exists": false,
		},
		concat(postAcceptedOffersKey, processEmail(vendorEmail)): types.M{
			"$exists": false,
		},
	}, options.Find().SetSort(types.M{
		updatedKey: 1,
	}).SetSkip(pageSize*pageNumber).SetLimit(pageSize).SetProjection(types.M{
		postOwnerKey:          0,
		postOffersKey:         0,
		postAcceptedOffersKey: 0,
	}))
}

// FetchOfferedPostsByVendor returns all posts the vendor has made an offer to
func FetchOfferedPostsByVendor(vendorEmail string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.OPEN,
		concat(postOffersKey, processEmail(vendorEmail)): types.M{
			"$exists": true,
		},
	}, options.Find().SetProjection(types.M{
		postOwnerKey: 0,
	}))
}

// FetchContractedPostsByVendor returns all posts in which the vendor's offer has been accepted
func FetchContractedPostsByVendor(vendorEmail string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.M{
			"$in": []string{types.OPEN, types.ONGOING},
		},
		concat(postAcceptedOffersKey, processEmail(vendorEmail)): types.M{
			"$exists": true,
		},
	}, options.Find().SetProjection(types.M{
		postOwnerKey: 0,
	}))
}

func fetchPostOffers(docID primitive.ObjectID, clientEmail string) (map[string]types.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err := postCollection.FindOne(ctx, types.M{
		primaryKey:   docID,
		postOwnerKey: clientEmail,
	}, options.FindOne().SetProjection(types.M{postOffersKey: 1})).Decode(post)
	return post.Offers, err
}

// FetchPostRequirements returns the requirements of a post
func FetchPostRequirements(postID string) (*types.Inventory, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postRequirementsKey: 1})).Decode(post)
	return &post.Requirements, err
}

// AcceptOffer accepts an offer made by a vendor on a post
// This operation is invoked by the client who is the owner of the post
// The param "offerKey" is key of the post holding the offer
// It is in the form of the vendor's email who made the offer with all "." replaced with "_"
// For Ex:- If the vendor's email is abc.2000@xyz.com the the key will be abc_2000@xyz_com
func AcceptOffer(postID, clientEmail, offerKey string) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey:   docID,
		postOwnerKey: clientEmail,
	}
	offers, err := fetchPostOffers(docID, clientEmail)
	if err != nil {
		return err
	}
	offerContent, ok := offers[offerKey]
	if !ok {
		return fmt.Errorf("Offer key %s doesnt exist in post %s for client %s", offerKey, postID, clientEmail)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return postCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$unset": types.M{
			concat(postOffersKey, offerKey): "",
		},
		"$set": types.M{
			concat(postAcceptedOffersKey, offerKey): offerContent,
		},
	}).Err()
}
