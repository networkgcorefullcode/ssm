package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// InsertData inserta un documento en MongoDB
func InsertData(client *mongo.Client, database string, collection string, data any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	result, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// InsertMultipleData inserta múltiples documentos en MongoDB
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
