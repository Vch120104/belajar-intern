package utils

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// convert Pascal to Snake Case for json response
func ModifyKeysInResponse(data interface{}) interface{} {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		newData := make(map[string]interface{})

		for _, key := range v.MapKeys() {
			newKey := PascaltoSnakeCase(key.String())
			value := v.MapIndex(key).Interface()
			newData[newKey] = ModifyKeysInResponse(value)
		}

		return newData
	case reflect.Slice:
		var newData []interface{}

		for i := 0; i < v.Len(); i++ {
			newData = append(newData, ModifyKeysInResponse(v.Index(i).Interface()))
		}

		return newData
	default:
		return data
	}
}

func PascaltoSnakeCase(input string) string {
	var output string
	for i, c := range input {
		if i > 0 && c >= 'A' && c <= 'Z' {
			output += "_"
		}
		output += string(c)
	}
	return strings.ToLower(output)
}

func SnaketoPascalCase(input string) string {
	words := strings.Split(input, "_")
	var pascalCase string
	for _, word := range words {
		if word != "" {
			pascalCase += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return pascalCase
}

// BuildFilterCondition insert columnValue and columnName into FilterCondition struct of slices then return it
func BuildFilterCondition(queryParams map[string]string) []FilterCondition {
	var criteria []FilterCondition
	for key, value := range queryParams {
		if value != "" {
			criteria = append(criteria, FilterCondition{
				ColumnValue: value,
				ColumnField: key,
			})
		}
	}
	return criteria
}

// Deprecated: please change to the latest one without *gin.Context
// GetQueryInt take QueryParam to return int
func GetQueryInt(c *gin.Context, param string) int {
	value, _ := strconv.Atoi(c.Query(param))
	return value
}

func NewGetQueryInt(queryValues url.Values, param string) int {
	value, _ := strconv.Atoi(queryValues.Get(param))
	return value
}

// ConvertDateTimeFormat converts ISO8601 format to separate date and time components
func ConvertDateFormat(dateTime time.Time) (string, error) {
	date := dateTime.Format("2006-01-02")
	return date, nil
}

// Converts a date string from the format '2 Jan 2006' to '2006-01-02'.
func ConvertDateStrFormat(dateStr string) (string, error) {
	parsedDate, err := time.Parse("2 Jan 2006", dateStr)
	if err != nil {
		return "", err
	}
	date := parsedDate.Format("2006-01-02")
	return date, nil
}

func ConvertTimeFormat(dateTimeStr string) (string, error) {
	// Parse the input string into time.Time
	parsedTime, err := time.Parse(time.RFC3339Nano, dateTimeStr)
	if err != nil {
		return "", err
	}

	// Extract time component
	time := parsedTime.Format("15:04:05")

	return time, nil
}

func ConvertDateTimeFormat(dateTimeStr string) (string, string, error) {
	// Parse the input string into time.Time
	parsedTime, err := time.Parse(time.RFC3339Nano, dateTimeStr)
	if err != nil {
		return "", "", err
	}

	// Extract date and time components
	date := parsedTime.Format("2006-01-02")
	time := parsedTime.Format("15:04:05")

	return date, time, nil
}

// SafeConvertDateFormat attempts to convert time.Time to string format "2006-01-02"
// If conversion fails, it returns an empty string
func SafeConvertDateFormat(dateTime time.Time) string {
	if dateTime.IsZero() {
		return ""
	}
	date, err := ConvertDateFormat(dateTime)
	if err != nil {
		return ""
	}
	return date
}

// Converts a date string from the format '2 Jan 2006' to '2006-01-02'.
// If conversion fails, returns an empty string.
func SafeConvertDateStrFormat(dateStr string) string {
	date, err := ConvertDateStrFormat(dateStr)
	if err != nil {
		return ""
	}
	return date
}

// Convert time.Time into time value. For example, converts '2024-01-01 15:34:45' into '15.579167'
func TimeValue(t time.Time) float64 {
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	timeValue := float64(hour) + (float64(minute)+float64(second)/60)/60

	return timeValue
}
