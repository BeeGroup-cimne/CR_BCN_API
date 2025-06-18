package utils

func ToStringArray(input interface{}) []string {
	interfaces, ok := input.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, len(interfaces))
	for i, v := range interfaces {
		if str, ok := v.(string); ok {
			result[i] = str
		}
	}
	return result
}
