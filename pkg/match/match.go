package match

import "regexp"

// EuroPlates returns whatever car license plate number is valid or not.
func EuroPlates(lexeme string) bool {
	matched, err := regexp.MatchString(`^\p{L}{2}\d{4}\p{L}{2}$`, lexeme)

	return err == nil && matched
}
