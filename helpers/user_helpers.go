package helpers

import (
	"context"

	"github.com/datarohit/go-restaurant-management-backend-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckDuplicateFields(ctx context.Context, email *string, phone *string) (bool, error) {
	if emailExists, err := RecordExists(ctx, userCollection, "email", email); err != nil {
		return false, err
	} else if emailExists {
		return true, nil
	}

	if phoneExists, err := RecordExists(ctx, userCollection, "phone", phone); err != nil {
		return false, err
	} else if phoneExists {
		return true, nil
	}

	return false, nil
}

func FindUserByEmail(ctx context.Context, email *string, user *models.User) error {
	return userCollection.FindOne(ctx, bson.M{"email": email}).Decode(user)
}

func FindUserByID(ctx context.Context, userId string, user *models.User) error {
	projection := bson.D{
		{Key: "password", Value: 0},
		{Key: "accessToken", Value: 0},
		{Key: "refreshToken", Value: 0},
	}
	return userCollection.FindOne(ctx, bson.M{"userId": userId}, &options.FindOneOptions{Projection: projection}).Decode(user)
}

func GetPaginatedUsers(ctx context.Context, skip int64, recordPerPage int64) ([]models.User, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: recordPerPage}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "password", Value: 0},
			{Key: "accessToken", Value: 0},
			{Key: "refreshToken", Value: 0},
		}}},
	}

	cursor, err := userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	totalCount, err := userCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
