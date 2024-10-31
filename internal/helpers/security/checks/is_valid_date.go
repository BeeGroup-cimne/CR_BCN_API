package checks

import (
	"time"
)

func IsValidDate(dateStr string) bool {
	_, err := time.Parse(time.RFC3339, dateStr)
	return err == nil
}
