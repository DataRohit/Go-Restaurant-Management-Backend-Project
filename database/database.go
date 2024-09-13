package database

import (
	"context"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/config"
	"github.com/datarohit/go-restaurant-management-backend-project/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func DBinstance() *mongo.Client {
	log := utils.GetLogger()

	mongoURI := config.GetEnv("MONGODB_URI", "mongodb://localhost:27017")
	connectTimeout := config.GetEnvAsInt("MONGODB_CONNECT_TIMEOUT", 10)

	log.Info("Connecting to MongoDB",
		zap.String("uri", mongoURI),
		zap.Int("timeout", connectTimeout))

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Error connecting to MongoDB", zap.Error(err))
		return nil
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("MongoDB ping failed", zap.Error(err))
		return nil
	}

	log.Info("Successfully connected to MongoDB")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	log := utils.GetLogger()

	databaseName := config.GetEnv("MONGODB_DATABASE", "restaurant")
	log.Info("Opening MongoDB collection",
		zap.String("database", databaseName),
		zap.String("collection", collectionName))

	collection := client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		log.Error("Failed to open MongoDB collection", zap.String("collection", collectionName))
		return nil
	}

	return collection
}
