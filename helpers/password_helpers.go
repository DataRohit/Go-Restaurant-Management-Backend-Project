package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password *string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword := string(bytes)
	return &hashedPassword
}

func VerifyPassword(hashedPassword, plainPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false, "login or password is incorrect"
	}
	return true, ""
}
