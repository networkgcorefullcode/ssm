package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// InsertData insert a mongoDB document
func InsertData(client *mongo.Client, database string, collection string, data any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DbContext.GenMutex.Lock()
	defer DbContext.GenMutex.Unlock()

	coll := client.Database(database).Collection(collection)
	result, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// InsertMultipleData insert multiplies mongoDB documents
func InsertMultipleData(client *mongo.Client, database string, collection string, data []any) ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	result, err := coll.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}
