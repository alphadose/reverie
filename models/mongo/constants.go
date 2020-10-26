package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// projectDatabase is the name of the database used for storing all of reverie's information
	projectDatabase = "reverie"

	// primaryKey is the primary key for mongoDB documents
	primaryKey = "_id"

	// PageSize is the number of documents retrieved in one pull for the vendor
	PageSize = 30

	// timeout is the context timeout for generic operations
	timeout = 5
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
