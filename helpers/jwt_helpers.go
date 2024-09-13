package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}

var (
	JWT_SECRET string = os.Getenv("JWT_SECRET")
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
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
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
