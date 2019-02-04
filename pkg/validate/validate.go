// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package validate

import "regexp"

func EURCarLicenseNumber(num string) bool {
	result, err := regexp.MatchString("[a-zA-Z]{2}[1-9]{4}[a-zA-Z]{2}", num)
	return err == nil && result
}
