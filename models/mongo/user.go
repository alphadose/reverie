package mongo

import (
	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertUser is an abstraction over UpdateOne which updates an application in mongoDB
// or inserts it if the corresponding document doesn't exist
func UpsertUser(filter types.M, data interface{}) error {
	return UpdateOne(UserCollection, filter, data, options.FindOneAndUpdate().SetUpsert(true))
}
