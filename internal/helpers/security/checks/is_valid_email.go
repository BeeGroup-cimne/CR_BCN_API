package checks

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9_.+\-]+@[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-.]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}
