// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

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
