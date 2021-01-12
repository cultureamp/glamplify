package env

import (
	"os"
	"strconv"
)

// GetString gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetInt gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetInt(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return i
}

// GetBool gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return b
}
