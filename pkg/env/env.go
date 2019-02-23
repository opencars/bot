package env

import (
	"fmt"
	"os"
)

// MustFetch panics if environment variable value is empty.
func MustFetch(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		panic(fmt.Sprintf("environment variable %s is missing or empty", name))
	}

	return value
}

// Fetch returns environment variable value or fallback in case of empty.
func Fetch(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}
