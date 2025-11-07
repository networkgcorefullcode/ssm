package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindAllData retrieves all documents from a collection
func FindAllData(client *mongo.Client, database string, collection string, filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// FindOneData retrieves a single document from a collection
func FindOneData(client *mongo.Client, database string, collection string, filter bson.M) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	var result bson.M
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByID retrieves a document by its ObjectID
func FindByID(client *mongo.Client, database string, collection string, id string) (bson.M, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	filter := bson.M{"_id": objectID}
	return FindOneData(client, database, collection, filter)
}

// FindWithOptions retrieves documents with pagination and sorting
func FindWithOptions(client *mongo.Client, database string, collection string, filter bson.M, opts *options.FindOptions) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// CountDocuments counts documents matching the filter
func CountDocuments(client *mongo.Client, database string, collection string, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// AggregateData performs aggregation operations
func AggregateData(client *mongo.Client, database string, collection string, pipeline []bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
