package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// OPEN denotes the status when a job request is open for offerings
	OPEN = "OPEN"
	// ONGOING denotes the status when a job request is in progress
	ONGOING = "ONGOING"
	// COMPLETED denotes the status when a job request is successfully completed
	COMPLETED = "COMPLETED"
	// DELETED denotes the status when a job request is removed by the client who posted it
	DELETED = "DELETED"
)

// Location denotes the location of the job request
type Location struct {
	// Always "Point"
	Type       string  `json:"-" bson:"type"`
	Latitude   float64 `json:"latitude" bson:"-"`
	Longtitude float64 `json:"longitude" bson:"-"`
	// Coordinates are in the form of [longitude, latitude] according to GeoJSON specifications
	Coordinates []float64 `json:"-" bson:"coordinates"`
	Place       string    `json:"place" bson:"place"`
}

// Post stores the information about a job request
type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Owner       string             `json:"owner" bson:"owner"`
	Description string             `json:"description" bson:"description"`
	Location    Location           `json:"location" bson:"location"`
	//  Infrastructure required by the client
	Requirements Inventory `json:"requirements" bson:"requirements"`
	// In the form of <email ID of the vendor offering the deal>:<the contents of the offer>
	Offers map[string]Inventory `json:"offers,omitempty" bson:"offers,omitempty"`
	// Status can be either OPEN, ONGOING, COMPLETED or DELETED
	Status  string `json:"-" bson:"status"`
	Created int64  `json:"-" bson:"created"`
	Updated int64  `json:"-" bson:"updated"`
}

// Initialize initializes the post parameters
func (post *Post) Initialize() {
	// Location
	post.Location.Coordinates = []float64{post.Location.Longtitude, post.Location.Latitude}
	post.Location.Type = "Point"

	// Timestamp
	post.Created = time.Now().Unix()
	post.Updated = time.Now().Unix()
}

// UpdateTimestamp updates the post's timestamp
func (post *Post) UpdateTimestamp() {
	post.Updated = time.Now().Unix()
}

// SetOwner sets the owner in the post's context
func (post *Post) SetOwner(ownerEmail string) {
	post.Owner = ownerEmail
}

// UpdateOffers updates the vendor's offerings in the post's context
func (post *Post) UpdateOffers(vendorEmail string, vendorOfferings Inventory) {
	post.Offers[vendorEmail] = vendorOfferings
}
