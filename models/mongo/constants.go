package mongo

import (
	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// projectDatabase is the name of the database used for storing all of reverie's information
	projectDatabase = "reverie"

	// PrimaryKey is the primary key for mongoDB documents
	PrimaryKey = "_id"

	// UserCollectionKey is the collection for all users
	UserCollectionKey = "users"

	// PostCollectionKey is the collection for all job posts
	PostCollectionKey = "posts"

	// NameKey is the key holding the name of an instance
	NameKey = "name"

	// EmailKey is the key holding the email of a user
	EmailKey = types.EmailKey

	// UsernameKey is the key holding the username of a user
	UsernameKey = types.UsernameKey

	// PasswordKey is the key holding the password of a user/instance
	PasswordKey = "password"

	// RoleKey is the key denoting the role of a user
	RoleKey = types.RoleKey

	// UserInventoryKey is the key denoting the inventory of a user
	UserInventoryKey = "inventory"

	// PostRequirementsKey is the key denoting the requirements for a post
	PostRequirementsKey = "requirements"

	// PostLocationKey is the key denoting the location of a job request
	PostLocationKey = "location"

	// CreatedKey is the key denoting the timestamp of creation of a job request
	CreatedKey = "created"

	// UpdatedKey is the key denoting the timestamp at which the job request was last updated
	UpdatedKey = "updated"

	// timeout is the context timeout for generic operations
	timeout = 5
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
