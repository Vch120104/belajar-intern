package utils

import (
	"after-sales/api/payloads/pagination"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

const CacheExpiration = 30 * time.Second // cache expiration time 30 seconds for live transaction

func GenerateCacheKeyIds(prefix string, params ...interface{}) string {

	var paramStrs []string
	for _, param := range params {
		switch v := param.(type) {
		case int:
			paramStrs = append(paramStrs, fmt.Sprintf("%d", v))
		case string:
			paramStrs = append(paramStrs, v)
		case []FilterCondition:
			filterBytes, _ := json.Marshal(v)
			paramStrs = append(paramStrs, string(filterBytes))
		case pagination.Pagination:
			paramStrs = append(paramStrs, fmt.Sprintf("page=%d&size=%d", v.Page, v.Limit))
		}
	}

	key := prefix + ":" + strings.Join(paramStrs, ":")

	return key
}

// Function to generate cache key for GetAll
func GenerateCacheKeys(prefix string, filterCondition []FilterCondition, pagination pagination.Pagination) string {

	filterBytes, _ := json.Marshal(filterCondition)

	pageStr := fmt.Sprintf("page=%d&size=%d", pagination.Page, pagination.Limit)

	key := fmt.Sprintf("%s:%s:%s", prefix, filterBytes, pageStr)

	return key

	// // pakai ini kalau ingin di hash key
	// hasher := sha1.New()
	// hasher.Write([]byte(key))
	// sha := hex.EncodeToString(hasher.Sum(nil))
	// return sha
}

// Function to refresh cache
func RefreshCaches(ctx context.Context, prefix interface{}) {
	var prefixStr string
	switch v := prefix.(type) {
	case string:
		prefixStr = v
	case int:
		prefixStr = strconv.Itoa(v)
	default:
		fmt.Println("Invalid prefix type. Must be string or int.")
		return
	}

	iter := RedisClient.Scan(ctx, 0, prefixStr+"*", 0).Iterator()
	for iter.Next(ctx) {
		RedisClient.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Error while scanning Redis keys:", err)
	}
}
