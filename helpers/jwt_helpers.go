package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/config"
	"github.com/datarohit/go-restaurant-management-backend-project/database"
	"github.com/datarohit/go-restaurant-management-backend-project/models"
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
	foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
	JWT_SECRET     string            = config.GetEnv("JWT_SECRET", "not-so-secret")
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
			ExpiresAt: time.Now().UTC().Add(6 * time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour).Unix(),
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

func GenerateToken(email, firstName, lastName, uid string, duration time.Duration) (string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UID:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(duration).Unix(),
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

func GetOrGenerateTokens(user models.User) (string, string, error) {
	if user.AccessToken != nil && ValidateToken(*user.AccessToken) == nil {
		return *user.AccessToken, *user.RefreshToken, nil
	}

	if user.RefreshToken != nil && ValidateToken(*user.RefreshToken) == nil {
		newAccessToken, err := GenerateToken(*user.Email, *user.FirstName, *user.LastName, user.UserID, 6*time.Hour)
		if err != nil {
			return "", "", err
		}
		return newAccessToken, *user.RefreshToken, nil
	}

	newAccessToken, newRefreshToken, err := GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
