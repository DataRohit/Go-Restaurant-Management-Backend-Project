package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetEnv(env, defaultValue string) string {
	environment := strings.TrimSpace(os.Getenv(env))
	if environment == "" {
		return defaultValue
	}
	return environment
}

func GetEnvAsInt(env string, defaultValue int) int {
	environment := strings.TrimSpace(os.Getenv(env))
	if environment == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(environment)
	if err != nil {
		log.Printf("Warning: %s is not a valid integer. Using default value: %d", env, defaultValue)
		return defaultValue
	}

	return value
}
