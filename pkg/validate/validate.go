package validate

import "regexp"

func EURCarLicenseNumber(num string) bool {
	result, err := regexp.MatchString("[a-zA-Z]{2}[1-9]{4}[a-zA-Z]{2}", num)
	return err == nil && result
}
