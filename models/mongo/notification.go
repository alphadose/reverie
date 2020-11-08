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

// FetchUnreadNotifications returns all unread notifications for a user
func FetchUnreadNotifications(email string) ([]types.M, error) {
	return fetchDocs(notificationCollection, types.M{
		notificationRecipentKey: email,
		notificationReadKey:     false,
	}, options.Find().SetSort(types.M{
		createdKey: -1,
	}))
}
