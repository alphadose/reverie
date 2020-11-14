package mongo

import (
	"fmt"
	"time"

	"github.com/reverie/types"
	"github.com/reverie/utils"
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

func notifyVendor(postID, vendorEmail, messageTemplate string) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	postName, err := FetchPostName(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	message := fmt.Sprintf(messageTemplate, postName)
	_, err = insertOne(notificationCollection, types.Notification{
		PostID:   docID,
		Recipent: vendorEmail,
		Type:     types.INFO,
		Message:  message,
		Read:     false,
		Created:  time.Now().Unix(),
	})
	if err != nil {
		utils.LogError("", err)
	}
}

// NotifyVendorOnAcceptance notifies a vendor when his offer on a post has been accepted
func NotifyVendorOnAcceptance(postID, vendorEmail string) {
	notifyVendor(postID, vendorEmail, "Your offer on post %s has been accepted")
}

// NotifyVendorOnRejection notifies a vendor when his offer on a post has been rejected
func NotifyVendorOnRejection(postID, vendorEmail string) {
	notifyVendor(postID, vendorEmail, "Your offer on post %s has been rejected")
}

// BulkNotifyVendors notfies all vendors whose offer has been accepted whenever there is a change in the post's status
func BulkNotifyVendors(postID, status string) {
	messageTemplate := ""

	switch status {
	case types.ONGOING:
		messageTemplate = "Work on post %s has started. Kindly deliver your equipments soon."
	case types.COMPLETED:
		messageTemplate = "Post %s has completed successfully"
	case types.DELETED:
		messageTemplate = "Post %s has been deleted"
	case types.OPEN:
		messageTemplate = "Work on post %s is temporarily halted"
	}

	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	acceptedOffers, postName, err := FetchPostAcceptedOffersAndName(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	message := fmt.Sprintf(messageTemplate, postName)

	offerKeys := make([]string, 0)
	for key := range acceptedOffers {
		offerKeys = append(offerKeys, key)
	}

	payload := make([]interface{}, 0)

	for _, offerKey := range offerKeys {
		vendorEmail, err := utils.Decrypt(offerKey)
		if err != nil {
			utils.LogError("kekw", err)
			continue
		}
		payload = append(payload, types.Notification{
			PostID:   docID,
			Recipent: vendorEmail,
			Type:     types.INFO,
			Message:  message,
			Read:     false,
			Created:  time.Now().Unix(),
		})
	}

	_, err = insertMany(notificationCollection, payload)
	if err != nil {
		utils.LogError("", err)
	}
}

// NotifyClient notifies a client whenever a vendor makes/retracts offer from his posts
func NotifyClient(postID, messageTemplate string) {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	postName, owner, err := FetchPostNameAndOwner(postID)
	if err != nil {
		utils.LogError("", err)
		return
	}
	message := fmt.Sprintf(messageTemplate, postName)
	_, err = insertOne(notificationCollection, types.Notification{
		PostID:   docID,
		Recipent: owner,
		Type:     types.INFO,
		Message:  message,
		Read:     false,
		Created:  time.Now().Unix(),
	})
	if err != nil {
		utils.LogError("", err)
	}
}

// NotifyOfferChangeToVendor notifies a vendor whenever a client requests changes on his offer
func NotifyOfferChangeToVendor(postID, vendorEmail string, offerChange *types.Inventory) error {
	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	postName, err := FetchPostName(postID)
	if err != nil {
		return err
	}
	_, err = insertOne(notificationCollection, types.Notification{
		PostID:         docID,
		Recipent:       vendorEmail,
		Type:           types.RequestOfferChange,
		Message:        fmt.Sprintf("Changes requested on post %s", postName),
		Read:           false,
		DesiredContent: offerChange,
		Created:        time.Now().Unix(),
	})
	return err
}
