package helper

import (
	"os"
	"strconv"
)

// GetEnvString gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetEnvString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvInt gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetEnvInt(key string, defaultValue int) int {
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

// GetEnvBool gets the environment variable for 'key' if present, otherwise returns 'fallback'
func GetEnvBool(key string, defaultValue bool) bool {
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
