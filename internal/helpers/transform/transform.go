package transform

import (
	"time"
)

func ConvertToNullableString(value interface{}) *string {
	if value == nil {
		return nil
	}
	str := value.(string)
	return &str
}

func ConvertToNullableInt(value interface{}) *int {
	if value == nil {
		return nil
	}
	integer := value.(int)
	return &integer
}

func ConvertToNullableBool(value interface{}) *bool {
	if value == nil {
		return nil
	}
	boolean := value.(bool)
	return &boolean
}

func ConvertIsoStringToTime(isoDate string) time.Time {
	layout := time.RFC3339
	parsedTime, err := time.Parse(layout, isoDate)
	if err != nil {
		return time.Now().UTC()
	}
	return parsedTime
}
