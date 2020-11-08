package types

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// REQUEST_OFFER_CHANGE is the type of notification sent when a client requests a vendor to change his offer
	REQUEST_OFFER_CHANGE = "REQUEST_OFFER_CHANGE"

	// INFO is the type of notification sent as means of informing clients/vendors for cases such as when an offer has been accepted, a post has started/completed etc
	INFO = "INFO"
)

// Notification is the schema for notifications
type Notification struct {
	// ID of the notification document in mongoDB
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// PostID is the ID of the post document in mongoDB with which the notification is concerned
	PostID primitive.ObjectID `json:"post_id" bson:"post_id"`

	// Type denotes the type of notification
	// Can be either REQUEST_OFFER_CHANGE or INFO
	Type string `json:"type" bson:"type"`

	// Message contains the body of the notification
	Message string `json:"message" bson:"message"`

	// Read denotes whether the notification is read or not
	Read bool `json:"read" bson:"read"`
}
