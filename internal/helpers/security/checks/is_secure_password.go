package checks

import (
	"regexp"
)

func IsSecurePassword(password string) bool {
	hasUppercase, _ := regexp.MatchString(`[A-Z]`, password)
	hasLowercase, _ := regexp.MatchString(`[a-z]`, password)
	hasNumber, _ := regexp.MatchString(`\d`, password)
	hasSymbol, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>_-]`, password)
	isLongEnough := len(password) >= 8

	return hasUppercase && hasLowercase && hasNumber && hasSymbol && isLongEnough
}
