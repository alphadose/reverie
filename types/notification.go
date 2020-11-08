package types

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// RequestOfferChange is the type of notification sent when a client requests a vendor to change his offer
	RequestOfferChange = "REQUEST_OFFER_CHANGE"

	// INFO is the type of notification sent as means of informing clients/vendors for cases such as when an offer has been accepted, a post has started/completed etc
	INFO = "INFO"
)

// TODO : cron job for cleaning up "Read" notifications older than 6 months

// Notification is the schema for notifications
type Notification struct {
	// ID of the notification document in mongoDB
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// PostID is the ID of the post document in mongoDB with which the notification is concerned
	PostID primitive.ObjectID `json:"post_id" bson:"post_id"`

	// Recipent is the email address of the recipent of the notification
	Recipent string `json:"recipent" bson:"recipent"`

	// Type denotes the type of notification
	// Can be either REQUEST_OFFER_CHANGE or INFO
	Type string `json:"type" bson:"type"`

	// Message contains the body of the notification
	Message string `json:"message" bson:"message"`

	// DesiredContent is the offer desired from a vendor by a client
	// If a client wants a vendor to change his offer, this field shall hold the final contents desired by the client
	// This field is only populated if the notification is of type REQUEST_OFFER_CHANGE
	DesiredContent *Inventory `json:"desired_content,omitempty" bson:"desired_content,omitempty"`

	// Read denotes whether the notification is read or not
	Read bool `json:"read" bson:"read"`

	Created int64 `json:"created" bson:"created"`
}
