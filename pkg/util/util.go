package util

import (
"fmt"
"os"
)

func MustGetEnv(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		panic(fmt.Sprintf("Environment variable %s does not set", name))
	}

	return value
}
