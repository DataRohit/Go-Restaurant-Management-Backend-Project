package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}

var (
	userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
	JWT_SECRET     string            = os.Getenv("JWT_SECRET")
)

func GenerateAllTokens(email, firstName, lastName, uid string) (string, string, error) {
	if JWT_SECRET == "" {
		return "", "", errors.New("JWT_SECRET is not set in the environment")
	}

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UID:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(6 * time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func UpdateAllTokens(signedAccessToken, signedRefreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.D{
		{Key: "accessToken", Value: signedAccessToken},
		{Key: "refreshToken", Value: signedRefreshToken},
		{Key: "updatedAt", Value: time.Now().UTC()},
	}

	filter := bson.M{"userId": userId}
	opts := options.Update().SetUpsert(true)

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateFields}}, opts)
	if err != nil {
		log.Printf("Failed to update tokens for user %s: %v", userId, err)
		return err
	}
	return nil
}

func ValidateToken(tokenStr string) error {
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	return err
}

func GenerateToken(email, uid string, duration time.Duration) (string, error) {
	claims := &SignedDetails{
		Email: email,
		UID:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
}

func UpdateAccessToken(signedAccessToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.D{
		{Key: "accessToken", Value: signedAccessToken},
		{Key: "updatedAt", Value: time.Now().UTC()},
	}

	filter := bson.M{"userId": userId}
	opts := options.Update().SetUpsert(true)

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateFields}}, opts)
	if err != nil {
		log.Printf("Failed to update access token for user %s: %v", userId, err)
		return err
	}
	return nil
}
