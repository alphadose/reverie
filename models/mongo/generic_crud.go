package mongo

import (
	"context"
	"time"

	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Operations

// insertOne inserts a document into a mongoDB collection
func insertOne(collection *mongo.Collection, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// insertMany inserts multiple document into a mongoDB collection
func insertMany(collection *mongo.Collection, data []interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// Read Operations

// fetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func fetchDocs(collection *mongo.Collection, filter types.M, opts ...*options.FindOptions) ([]types.M, error) {
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

// countDocs returns the number of documents matching a filter
func countDocs(collection *mongo.Collection, filter types.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}

// Update Operations

// updateOne updates a document in the mongoDB collection
func updateOne(collection *mongo.Collection, filter types.M, data interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.FindOneAndUpdate(ctx, filter, types.M{"$set": data}, opts...).Err()
}

// bulkUpsert upserts multiple documents using BulkWrite
func bulkUpsert(collection *mongo.Collection, data []mongo.WriteModel, opts ...*options.BulkWriteOptions) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.BulkWrite(ctx, data, opts...)
}

// updateMany updates multiple documents in the mongoDB collection
func updateMany(collection *mongo.Collection, filter types.M, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.UpdateMany(ctx, filter, types.M{"$set": data}, nil)
}

// Delete Operations

// deleteOne deletes a document from a mongoDB collection
func deleteOne(collection *mongo.Collection, filter types.M) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}
