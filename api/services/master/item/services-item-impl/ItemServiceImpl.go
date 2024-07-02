package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const cacheExpiration = time.Minute * 30 // cache expiration time

type ItemServiceImpl struct {
	itemRepo    masteritemrepository.ItemRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartItemService(itemRepo masteritemrepository.ItemRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemService {
	return &ItemServiceImpl{
		itemRepo:    itemRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

// Function to generate cache key for GetAll
func generateCacheKey(prefix string, filterCondition []utils.FilterCondition, pagination pagination.Pagination) string {

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

// Function to generate cache key for GetById
func generateCacheKeyId(prefix string, params ...interface{}) string {

	var paramStrs []string
	for _, param := range params {
		switch v := param.(type) {
		case int:
			paramStrs = append(paramStrs, fmt.Sprintf("%d", v))
		case string:
			paramStrs = append(paramStrs, v)
		case []utils.FilterCondition:
			filterBytes, _ := json.Marshal(v)
			paramStrs = append(paramStrs, string(filterBytes))
		case pagination.Pagination:
			paramStrs = append(paramStrs, fmt.Sprintf("page=%d&size=%d", v.Page, v.Limit))
		}
	}

	key := prefix + ":" + strings.Join(paramStrs, ":")

	return key
}

// GetUomDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomDropDown(tx, uomTypeId)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

// GetUomTypeDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomTypeDropDown(tx)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("item_master", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("Cache miss, querying database...")

		tx := s.DB.Begin()

		results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItem(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr != nil {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		} else {
			setErr := s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration).Err()
			if setErr != nil {
				fmt.Println("Failed to set cache:", setErr)
			}
		}
		defer helper.CommitOrRollback(tx, repoErr)
		return results, totalPages, totalRows, nil
	} else if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	} else {
		fmt.Println("Cache hit, returning cached data...")
		var mapResponses []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &mapResponses); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
		return paginatedData, totalPages, totalRows, nil
	}
}

func (s *ItemServiceImpl) GetAllItemLookup(filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetAllItemLookup(tx, filter)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	// Melakukan URL encoding pada parameter code
	// encodedCode := url.PathEscape(code)

	tx := s.DB.Begin()
	results, err := s.itemRepo.GetItemCode(tx, code) // Menggunakan kode yang telah diencode
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	fmt.Print("sini?")

	if req.ItemId != 0 {
		_, err := s.itemRepo.GetItemById(tx, req.ItemId)
		if err != nil {
			return false, err

		}
	}

	results, err := s.itemRepo.SaveItem(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("item_detail_master", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("Cache miss, querying database...")
		tx := s.DB.Begin()
		results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItemDetail(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr != nil {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		} else {
			setErr := s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration).Err()
			if setErr != nil {
				fmt.Println("Failed to set cache:", setErr)
			}
		}
		defer helper.CommitOrRollback(tx, repoErr)
		return results, totalPages, totalRows, nil
	} else if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	} else {
		fmt.Println("Cache hit, returning cached data...")
		var mapResponses []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &mapResponses); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
		return paginatedData, totalPages, totalRows, nil
	}
}

func (s *ItemServiceImpl) GetItemDetailById(itemID, itemDetailID int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKeyId("item_detail", itemID, itemDetailID)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var result masteritempayloads.ItemDetailRequest
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &result); unmarshalErr != nil {
			return masteritempayloads.ItemDetailRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        unmarshalErr,
			}
		}
		return result, nil
	} else if err != redis.Nil {

		return masteritempayloads.ItemDetailRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	result, repoErr := s.itemRepo.GetItemDetailById(tx, itemID, itemDetailID)
	if repoErr != nil {
		errorResponse := &exceptions.BaseErrorResponse{Message: repoErr.Message}
		return masteritempayloads.ItemDetailRequest{}, errorResponse
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal results for caching:", marshalErr)
	} else {
		setErr := s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration).Err()
		if setErr != nil {
			fmt.Println("Failed to set cache:", setErr)
		}
	}
	defer helper.CommitOrRollback(tx, repoErr)
	return result, nil
}

func (s *ItemServiceImpl) AddItemDetail(id int, req masteritempayloads.ItemDetailRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.itemRepo.AddItemDetail(tx, id, req)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx, err)
	return nil
}

func (s *ItemServiceImpl) DeleteItemDetail(id int, itemDetailID int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.itemRepo.DeleteItemDetail(tx, id, itemDetailID)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx, err)
	return nil
}

func (s *ItemServiceImpl) UpdateItem(id int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.UpdateItem(tx, id, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) UpdateItemDetail(id int, req masteritempayloads.ItemDetailUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.UpdateItemDetail(tx, id, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandParent(code string) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandParent(tx, code)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandDropdown() ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandDropdown(tx)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}
