// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package env

import (
	"fmt"
	"os"
)

// MustGetEnv panics if environment variable does not exist, otherwise returns value.
func MustGet(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		panic(fmt.Sprintf("Environment variable %s does not set", name))
	}

	return value
}

// GetEnv returns environment variable with ability to specify default value.
func Get(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}
