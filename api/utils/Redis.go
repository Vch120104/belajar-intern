package utils

import "fmt"

func GenerateCacheKey(filterCondition []FilterCondition) string {
	return "cache_key_" + fmt.Sprint(filterCondition)
}
