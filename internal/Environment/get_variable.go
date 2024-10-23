package Environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	return os.Getenv(key)
}

func GetVars(keys []string) map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	vars := make(map[string]string)

	for _, key := range keys {
		value := os.Getenv(key)
		vars[key] = value
	}

	return vars
}
