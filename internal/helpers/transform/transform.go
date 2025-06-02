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

func ConvertNullableToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func ConvertToNullableInt(value interface{}) *int {
	if value == nil {
		return nil
	}
	var integer int
	switch v := value.(type) {
	case int:
		integer = v
	case int64:
		integer = int(v)
	case int32:
		integer = int(v)
	case int16:
		integer = int(v)
	case int8:
		integer = int(v)
	case uint:
		integer = int(v)
	case uint64:
		integer = int(v)
	case uint32:
		integer = int(v)
	case uint16:
		integer = int(v)
	case uint8:
		integer = int(v)
	default:
		return nil
	}

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

func ConvertToNullableTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case time.Time:
		return &t
	case *time.Time:
		return t
	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			return &parsed
		}
	}

	return nil
}
