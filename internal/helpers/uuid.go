package helpers

import (
	"regexp"
)

func IsValidULID(id string) bool {
	ulidRegex := regexp.MustCompile(`^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`)
	return ulidRegex.MatchString(id) && len(id) == 26
}
