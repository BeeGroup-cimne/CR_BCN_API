package env

import (
	"os"
)

func GetVar(key string) string {
	return os.Getenv(key)
}

func GetVars(keys []string) map[string]string {
	vars := make(map[string]string)
	for _, key := range keys {
		value := os.Getenv(key)
		vars[key] = value
	}
	return vars
}
