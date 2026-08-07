package cassandra

import "strings"

// propertyTranslation maps the English property names used in the public
// API URLs to the actual Spanish terms stored in the Cassandra hash
// (e.g. https://climres.cat#skyspark-measurement-am44-{sensorID}-{dwelling}-temperatura-PT15M).
var propertyTranslation = map[string]string{
	"temperature": "temperatura",
	"humidity":    "humedad",
	"co2":         "co2",
	"battery":     "bateria",
}

// TranslateProperty converts a public (English) property name into the
// term actually stored in Cassandra. If the property is not recognized,
// it's returned unchanged (lowercased) so callers can still try it as-is.
func TranslateProperty(property string) string {
	p := strings.ToLower(property)
	if translated, ok := propertyTranslation[p]; ok {
		return translated
	}
	return p
}
