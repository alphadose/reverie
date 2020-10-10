package mongo

import (
	"context"
	"time"

	"github.com/reverie/types"
	"github.com/reverie/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne inserts a document into a mongoDB collection
func InsertOne(collectionName string, data interface{}) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// InsertMany inserts multiple document into a mongoDB collection
func InsertMany(collectionName string, data []interface{}) ([]interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collectionName string, filter types.M, opts ...*options.FindOptions) []types.M {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var data []types.M

	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		utils.LogError("Mongo-Read-1", err)
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result types.M
		if err := cur.Decode(&result); err != nil {
			utils.LogError("Mongo-Read-2", err)
			return nil
		}
		data = append(data, result)
	}
	if err := cur.Err(); err != nil {
		utils.LogError("Mongo-Read-3", err)
		return nil
	}
	return data
}

// UpdateOne updates a document in the mongoDB collection
func UpdateOne(collectionName string, filter types.M, data interface{}, option *options.FindOneAndUpdateOptions) error {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.FindOneAndUpdate(ctx, filter, types.M{"$set": data}, option).Err()
}

// BulkUpsert upserts multiple documents using BulkWrite
func BulkUpsert(collectionName string, data []mongo.WriteModel, options *options.BulkWriteOptions) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.BulkWrite(ctx, data, options)
}

// UpdateMany updates multiple documents in the mongoDB collection
func UpdateMany(collectionName string, filter types.M, data interface{}) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.UpdateMany(ctx, filter, types.M{"$set": data}, nil)
}

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collectionName string, filter types.M) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}
