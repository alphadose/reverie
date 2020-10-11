package mongo

import (
	"context"
	"time"

	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Operations

// InsertOne inserts a document into a mongoDB collection
func InsertOne(collection *mongo.Collection, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// InsertMany inserts multiple document into a mongoDB collection
func InsertMany(collection *mongo.Collection, data []interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// Read Operations

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collection *mongo.Collection, filter types.M, opts ...*options.FindOptions) ([]types.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	var data []types.M
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &data)
	return data, err
}

// CountDocs returns the number of documents matching a filter
func CountDocs(collection *mongo.Collection, filter types.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}

// Update Operations

// UpdateOne updates a document in the mongoDB collection
func UpdateOne(collection *mongo.Collection, filter types.M, data interface{}, option *options.FindOneAndUpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.FindOneAndUpdate(ctx, filter, types.M{"$set": data}, option).Err()
}

// BulkUpsert upserts multiple documents using BulkWrite
func BulkUpsert(collection *mongo.Collection, data []mongo.WriteModel, options *options.BulkWriteOptions) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.BulkWrite(ctx, data, options)
}

// UpdateMany updates multiple documents in the mongoDB collection
func UpdateMany(collection *mongo.Collection, filter types.M, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.UpdateMany(ctx, filter, types.M{"$set": data}, nil)
}

// Delete Operations

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collection *mongo.Collection, filter types.M) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}
