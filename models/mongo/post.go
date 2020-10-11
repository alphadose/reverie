package mongo

import (
	"fmt"

	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postCollection = db.Collection(PostCollectionKey)

// CreatePost is an abstraction over InsertOne which inserts a post
func CreatePost(post *types.Post) (interface{}, error) {
	return InsertOne(postCollection, post)
}

// UpdatePostOffers adds/updates an offer to a post
func UpdatePostOffers(postID, vendorEmail string, offer *types.Inventory) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	filter := types.M{
		PrimaryKey: docID,
	}
	updatePayload := types.M{
		fmt.Sprintf("%s.%s", PostRequirementsKey, vendorEmail): offer,
	}
	return UpdateOne(postCollection, filter, updatePayload, options.FindOneAndUpdate().SetUpsert(true))
}
