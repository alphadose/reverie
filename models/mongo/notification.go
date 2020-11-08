package mongo

import (
	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// notificationCollectionKey is the collection for all notifications
	notificationCollectionKey = "notifications"

	// notificationRecipentKey is the key denoting the recipent (i.e email address) of the notification
	notificationRecipentKey = "recipent"

	// notificationReadKey denotes if a notification is read or not
	notificationReadKey = "read"

	// notificationPageSize is the maximum number of notifications per batch
	notificationPageSize = 20
)

var notificationCollection = db.Collection(notificationCollectionKey)

// MarkRead marks a notification as read by the user
func MarkRead(notificationID, email string) error {
	docID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey:              docID,
		notificationRecipentKey: email,
	}
	updatePayload := types.M{
		notificationReadKey: true,
	}
	return updateOne(notificationCollection, filter, updatePayload)
}

// FetchNotifications returns all notifications for a user
func FetchNotifications(email string, pageNumber int64) ([]types.M, error) {
	return fetchDocs(notificationCollection, types.M{
		notificationRecipentKey: email,
	}, options.Find().SetSort(types.M{
		createdKey: -1,
	}).SetSkip(notificationPageSize*pageNumber).SetLimit(notificationPageSize))
}
