package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// projectDatabase is the name of the database used for storing all of reverie's information
	projectDatabase = "reverie"

	// UserCollection is the collection for all users
	UserCollection = "users"

	// PostCollection is the collection for all job posts
	PostCollection = "posts"

	// NameKey is the key holding the name of an instance
	NameKey = "name"

	// EmailKey is the key holding the email of a user
	EmailKey = "email"

	// UsernameKey is the key holding the username of a user
	UsernameKey = "username"

	// PasswordKey is the key holding the password of a user/instance
	PasswordKey = "password"

	// RoleKey is the key denoting the role of a user
	RoleKey = "role"

	// timeout is the context timeout for generic operations
	timeout = 5
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
