package utils

import (
	"go-api/internal/models"
	"go-api/internal/models/enums/timePrecision"
	"time"
)

func IsEqualDate(t time.Time, d time.Time, p models.TimePrecision) bool {
	switch p {
	case timePrecision.Year:
		return t.Year() == d.Year()
	case timePrecision.Month:
		return t.Year() == d.Year() && t.Month() == d.Month()
	case timePrecision.Day:
		return t.Year() == d.Year() && t.Month() == d.Month() && t.Day() == d.Day()
	case timePrecision.Hour:
		return t.Year() == d.Year() && t.Month() == d.Month() && t.Day() == d.Day() && t.Hour() == d.Hour()
	case timePrecision.Minute:
		return t.Year() == d.Year() && t.Month() == d.Month() && t.Day() == d.Day() && t.Hour() == d.Hour() && t.Minute() == d.Minute()
	case timePrecision.Second:
		return t.Year() == d.Year() && t.Month() == d.Month() && t.Day() == d.Day() && t.Hour() == d.Hour() && t.Minute() == d.Minute() && t.Second() == d.Second()
	default:
		return t.Equal(d)
	}
}
