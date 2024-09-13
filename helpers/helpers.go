package helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RecordExists(ctx context.Context, collection *mongo.Collection, field string, value interface{}) (bool, error) {
	count, err := collection.CountDocuments(ctx, bson.M{field: value})
	return count > 0, err
}

func GetNonNilString(s *string, defaultValue string) string {
	if s == nil || *s == "" {
		return defaultValue
	}
	return *s
}
