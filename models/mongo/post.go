package mongo

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/reverie/types"
	"github.com/reverie/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// postCollectionKey is the collection for all job posts
	postCollectionKey = "posts"

	// postRequirementsKey is the key denoting the name of a post
	postNameKey = "name"

	// postRequirementsKey is the key denoting the description of a post
	postDescriptionKey = "description"

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

	// postOwnerNameKey is the key holding the owner's name of a post
	postOwnerNameKey = "owner_name"

	// postStatusKey is the key holding the status of a post
	postStatusKey = "status"

	// createdKey is the key denoting the timestamp of creation of a job request
	createdKey = "created"

	// updatedKey is the key denoting the timestamp at which the job request was last updated
	updatedKey = "updated"

	// postPageSize is the maximum number of posts retrieved in one batch for the vendor
	postPageSize = 30
)

// Constants for offer schema
const (
	// the name of the vendor making the offer
	offerNameKey = "name"

	// contents of the offer in the form of types.Inventory
	offerContentKey = "content"

	// time of creation of the offer
	offerTimestampKey = "created"
)

var postCollection = db.Collection(postCollectionKey)

// concatenates strings with "." in between
func concat(keys ...string) string {
	return strings.Join(keys, ".")
}

// CreatePost is an abstraction over InsertOne which inserts a post
func CreatePost(post *types.Post) (interface{}, error) {
	return insertOne(postCollection, post)
}

// IsPostOwner checks if a client is the owner of a post or not
func IsPostOwner(postID, clientEmail string) (bool, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, err
	}
	count, err := countDocs(postCollection, types.M{
		postOwnerKey: clientEmail,
		primaryKey:   docID,
	})
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

// UpdatePost updates a post by a client
func UpdatePost(postID string, post *types.PostUpdate) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	return updateOne(postCollection, filter, post)
}

// UpdatePostOffers adds/updates an offer to an OPEN post
func UpdatePostOffers(postID, vendorEmail string, offer types.Offer) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return err
	}
	updatePayload := types.M{
		concat(postOffersKey, vendorEmailKey): offer,
	}
	return updateOne(postCollection, filter, updatePayload)
}

// RetractPostOffer removes an offer from an OPEN post by a vendor
func RetractPostOffer(postID, vendorEmail string) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return postCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$unset": types.M{
			concat(postOffersKey, vendorEmailKey): "",
		},
	}).Err()
}

// RejectAcceptedOffer removes an accepted offer from an OPEN post by a client
// The param "offerKey" is key holding the offer in the post
// It is the vendor's email address encrypted with AES-256
func RejectAcceptedOffer(postID, offerKey string, offer types.Inventory) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}

	// Make a map for incrementing the posts's current requirements
	incrementMap := make(map[string]int64)

	offerValues := reflect.ValueOf(offer)
	offerKeys := reflect.TypeOf(offer)

	for i := 0; i < offerValues.NumField(); i++ {
		value := offerValues.Field(i).Int()
		if value == 0 {
			continue
		}
		key := concat(postRequirementsKey, offerKeys.Field(i).Name)
		incrementMap[key] = value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return postCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$unset": types.M{
			concat(postAcceptedOffersKey, offerKey): "",
		},
		"$inc": incrementMap,
	}).Err()
}

// RejectPendingOffer removes a pending offer from an OPEN post by a client
// The param "offerKey" is key holding the offer in the post
// It is the vendor's email address encrypted with AES-256
func RejectPendingOffer(postID, offerKey string) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return postCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$unset": types.M{
			concat(postOffersKey, offerKey): "",
		},
	}).Err()
}

// FetchActivePostsByClient returns all open/ongoing posts created by a client
func FetchActivePostsByClient(clientEmail string) ([]types.M, error) {
	return fetchDocs(postCollection, types.M{
		postOwnerKey: clientEmail,
		postStatusKey: types.M{
			"$in": []string{types.OPEN, types.ONGOING},
		},
	}, options.Find().SetSort(types.M{
		updatedKey: -1,
	}))
}

// FetchSinglePostByClient returns a single post given its id
func FetchSinglePostByClient(postID string) (*types.Post, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne()).Decode(post)

	return post, err
}

// UpdatePostStatus updates the status of the post
func UpdatePostStatus(postID, newStatus string) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	updatePayload := types.M{
		postStatusKey: newStatus,
	}
	return updateOne(postCollection, filter, updatePayload)
}

// FetchSinglePostByVendor returns a single post given its id
func FetchSinglePostByVendor(postID, vendorEmail string) (*types.Post, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, err
	}
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{ // TODO : update these fields as more information is added to posts
		postNameKey:                           1,
		postDescriptionKey:                    1,
		postOwnerNameKey:                      1,
		postLocationKey:                       1,
		postRequirementsKey:                   1,
		createdKey:                            1,
		concat(postOffersKey, vendorEmailKey): 1,
		concat(postAcceptedOffersKey, vendorEmailKey): 1,
	})).Decode(post)

	return post, err
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
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return []types.M{}, err
	}
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.OPEN,
		"$or":         searchArray,
		concat(postOffersKey, vendorEmailKey): types.M{
			"$exists": false,
		},
		concat(postAcceptedOffersKey, vendorEmailKey): types.M{
			"$exists": false,
		},
	}, options.Find().SetSort(types.M{
		updatedKey: -1,
	}).SetSkip(postPageSize*pageNumber).SetLimit(postPageSize).SetProjection(types.M{
		postOwnerKey:          0,
		postOffersKey:         0,
		postAcceptedOffersKey: 0,
	}))
}

// FetchOfferedPostsByVendor returns all posts the vendor has made an offer to
func FetchOfferedPostsByVendor(vendorEmail string) ([]types.M, error) {
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return []types.M{}, err
	}
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.OPEN,
		concat(postOffersKey, vendorEmailKey): types.M{
			"$exists": true,
		},
	}, options.Find().SetSort(types.M{
		updatedKey: -1,
	}).SetProjection(types.M{ // TODO : update these fields as more information is added to posts
		postNameKey:                           1,
		postDescriptionKey:                    1,
		postLocationKey:                       1,
		postRequirementsKey:                   1,
		postOwnerNameKey:                      1,
		createdKey:                            1,
		concat(postOffersKey, vendorEmailKey): 1,
	}))
}

// FetchContractedPostsByVendor returns all posts in which the vendor's offer has been accepted
func FetchContractedPostsByVendor(vendorEmail string) ([]types.M, error) {
	vendorEmailKey, err := utils.Encrypt(vendorEmail)
	if err != nil {
		return []types.M{}, err
	}
	return fetchDocs(postCollection, types.M{
		postStatusKey: types.M{
			"$in": []string{types.OPEN, types.ONGOING},
		},
		concat(postAcceptedOffersKey, vendorEmailKey): types.M{
			"$exists": true,
		},
	}, options.Find().SetSort(types.M{
		updatedKey: -1,
	}).SetProjection(types.M{ // TODO : update these fields as more information is added to posts
		postNameKey:         1,
		postDescriptionKey:  1,
		postLocationKey:     1,
		postRequirementsKey: 1,
		postOwnerNameKey:    1,
		createdKey:          1,
		concat(postAcceptedOffersKey, vendorEmailKey): 1,
	}))
}

// FetchPostStatus returns a post's status
func FetchPostStatus(postID string) (string, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	postStatus := &types.PostStatus{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postStatusKey: 1})).Decode(postStatus)
	if err != nil {
		return "", err
	}

	return postStatus.Value, nil
}

// FetchPostName returns a post's name
func FetchPostName(postID string) (string, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	name := &types.PostName{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postNameKey: 1})).Decode(name)
	if err != nil {
		return "", err
	}

	return name.Value, nil
}

// FetchPostNameAndOwner returns a post's name and owner
func FetchPostNameAndOwner(postID string) (string, string, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	contract := &types.PostNameAndOwner{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postNameKey: 1, postOwnerKey: 1})).Decode(contract)
	if err != nil {
		return "", "", err
	}

	return contract.Name, contract.Owner, nil
}

// FetchPostOffersAndRequirementsAndStatus returns a post's offers (both accepted and pending) and requirements as well as its status
func FetchPostOffersAndRequirementsAndStatus(postID string) (string, map[string]types.Offer, types.Inventory, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", nil, types.Inventory{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postOffersKey: 1, postRequirementsKey: 1, postStatusKey: 1})).Decode(post)
	if err != nil {
		return "", nil, types.Inventory{}, err
	}
	return post.Status, post.Offers, post.Requirements, nil
}

// FetchPostRequirementsAndStatus returns the requirements of a post
func FetchPostRequirementsAndStatus(postID string) (string, *types.Inventory, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return "", nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postRequirementsKey: 1, postStatusKey: 1})).Decode(post)
	if err != nil {
		return "", nil, err
	}
	return post.Status, &post.Requirements, nil
}

// FetchPostAcceptedOffersAndStatus returns the accepted offers of a post as well as its status
func FetchPostAcceptedOffersAndStatus(postID string) (map[string]types.Offer, string, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postAcceptedOffersKey: 1, postStatusKey: 1})).Decode(post)
	if err != nil {
		return nil, "", err
	}
	return post.AcceptedOffers, post.Status, nil
}

// FetchPostAcceptedOffersAndName returns the accepted offers of a post as well as its name
func FetchPostAcceptedOffersAndName(postID string) (map[string]types.Offer, string, error) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	post := &types.Post{}
	err = postCollection.FindOne(ctx, types.M{
		primaryKey: docID,
	}, options.FindOne().SetProjection(types.M{postAcceptedOffersKey: 1, postNameKey: 1})).Decode(post)
	if err != nil {
		return nil, "", err
	}
	return post.AcceptedOffers, post.Name, nil
}

// AcceptOffer accepts an offer made by a vendor on a post
// This operation is invoked by the client who is the owner of the post
// The param "offerKey" is key holding the offer in the post
// It is the vendor's email address encrypted with AES-256
func AcceptOffer(postID, offerKey string, offer types.Offer) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	filter := types.M{
		primaryKey: docID,
	}

	// Make a map for incrementing the post's accepted offers and decrementing the posts's current requirements
	// This is handy if a vendor makes offer twice and both are accepted
	// This section would combine the 2 individual offers into a single accepted offer
	changeMap := make(map[string]int64)

	offerValues := reflect.ValueOf(offer.Content)
	offerKeys := reflect.TypeOf(offer.Content)

	for i := 0; i < offerValues.NumField(); i++ {
		value := offerValues.Field(i).Int()
		if value == 0 {
			continue
		}
		acceptedOfferIncrementKey := concat(postAcceptedOffersKey, offerKey, offerContentKey, offerKeys.Field(i).Name)
		requirementsDecrementKey := concat(postRequirementsKey, offerKeys.Field(i).Name)
		changeMap[acceptedOfferIncrementKey] = value
		changeMap[requirementsDecrementKey] = value * -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return postCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$unset": types.M{
			concat(postOffersKey, offerKey): "",
		},
		"$inc": changeMap,
		// TODO : Update $set as more and more fields are added to offer schema
		"$set": types.M{
			concat(postAcceptedOffersKey, offerKey, offerNameKey):      offer.Name,
			concat(postAcceptedOffersKey, offerKey, offerTimestampKey): offer.Created,
		},
	}).Err()
}
