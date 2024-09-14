package helpers

import (
	"context"

	"github.com/datarohit/go-restaurant-management-backend-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPaginatedFoodItems(ctx context.Context, skip int64, recordPerPage int64) ([]models.Food, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: recordPerPage}},
	}

	cursor, err := foodCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var foodItems []models.Food
	if err := cursor.All(ctx, &foodItems); err != nil {
		return nil, 0, err
	}

	totalCount, err := foodCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return foodItems, totalCount, nil
}
