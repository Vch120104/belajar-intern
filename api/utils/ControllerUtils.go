package utils

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

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


