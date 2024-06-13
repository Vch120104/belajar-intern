package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const cacheExpiration = time.Minute * 5 // cache expiration time

type WorkOrderServiceImpl struct {
	structWorkOrderRepo transactionworkshoprepository.WorkOrderRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func OpenWorkOrderServiceImpl(WorkOrderRepo transactionworkshoprepository.WorkOrderRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.WorkOrderService {
	return &WorkOrderServiceImpl{
		structWorkOrderRepo: WorkOrderRepo,
		DB:                  db,
		RedisClient:         redisClient,
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

// Function to refresh cache
func (s *WorkOrderServiceImpl) refreshCache(ctx context.Context, prefix interface{}) {
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

	iter := s.RedisClient.Scan(ctx, 0, prefixStr+"*", 0).Iterator()
	for iter.Next(ctx) {
		s.RedisClient.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Error while scanning Redis keys:", err)
	}
}

func (s *WorkOrderServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("work_orders", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {

		fmt.Println("Cache miss, querying database...")

		tx := s.DB.Begin()
		defer helper.CommitOrRollback(tx)

		results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAll(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr == nil {
			s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
		} else {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		}

		// Refresh cache
		s.refreshCache(ctx, "work_orders")

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

func (s *WorkOrderServiceImpl) VehicleLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("vehicle_lookup", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {

		fmt.Println("Cache miss for VehicleLookup, querying database...")

		tx := s.DB.Begin()
		defer helper.CommitOrRollback(tx)

		results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.VehicleLookup(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr == nil {
			s.RedisClient.Set(ctx, cacheKey, cacheData, 5*time.Minute) // Atur durasi cache sesuai kebutuhan
		} else {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		}

		return results, totalPages, totalRows, nil
	} else if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	} else {

		fmt.Println("Cache hit for VehicleLookup, returning cached data...")
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

func (s *WorkOrderServiceImpl) CampaignLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("campaign_lookup", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {

		fmt.Println("Cache miss for CampaignLookup, querying database...")

		tx := s.DB.Begin()
		defer helper.CommitOrRollback(tx)

		results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.CampaignLookup(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr == nil {
			s.RedisClient.Set(ctx, cacheKey, cacheData, 5*time.Minute) // Atur durasi cache sesuai kebutuhan
		} else {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		}

		return results, totalPages, totalRows, nil
	} else if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	} else {

		fmt.Println("Cache hit for CampaignLookup, returning cached data...")
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

func (s *WorkOrderServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) (bool, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.New(tx, request)
	if err != nil {
		return false, err
	}

	// Refresh cache after adding new data
	s.refreshCache(ctx, "work_orders")

	return save, nil
}

func (s *WorkOrderServiceImpl) NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	statuses, err := s.structWorkOrderRepo.NewStatus(tx, filter)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *WorkOrderServiceImpl) NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	types, err := s.structWorkOrderRepo.NewType(tx, filter)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (s *WorkOrderServiceImpl) NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	bills, err := s.structWorkOrderRepo.NewBill(tx)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) NewDropPoint(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
	dropPoints, err := s.structWorkOrderRepo.NewDropPoint(tx)
	if err != nil {
		return nil, err
	}
	return dropPoints, nil
}

func (s *WorkOrderServiceImpl) NewVehicleBrand(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
	brands, err := s.structWorkOrderRepo.NewVehicleBrand(tx)
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (s *WorkOrderServiceImpl) NewVehicleModel(tx *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
	models, err := s.structWorkOrderRepo.NewVehicleModel(tx, brandId)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *WorkOrderServiceImpl) GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	idString := strconv.Itoa(id)
	cacheKey := generateCacheKeyId(idString)

	// retrieve data from cache
	cachedData, err := s.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {

		var result transactionworkshoppayloads.WorkOrderRequest
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Error unmarshalling cached data"}
		}
		return result, nil
	}

	// Data not found in cache, proceed to database
	// Start database transaction
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Retrieve data from repository
	results, repoErr := s.structWorkOrderRepo.GetById(tx, id)
	if repoErr != nil {

		errorResponse := &exceptions.BaseErrorResponse{Message: repoErr.Message}
		return transactionworkshoppayloads.WorkOrderRequest{}, errorResponse
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Error marshalling data"}
	}

	if err := s.RedisClient.Set(context.Background(), cacheKey, jsonData, cacheExpiration).Err(); err != nil {
		fmt.Println("Error caching data:", err)
	}

	// Refresh cache after adding new data
	s.refreshCache(ctx, idString)

	return results, nil
}

func (s *WorkOrderServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	// Start a new transaction
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.Save(tx, request, workOrderId)
	if err != nil {
		return false, err
	}

	// Refresh cache after adding new data
	s.refreshCache(ctx, workOrderId)

	return save, nil
}

func (s *WorkOrderServiceImpl) Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	delete, err := s.structWorkOrderRepo.Void(tx, workOrderId)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) CloseOrder(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	close, err := s.structWorkOrderRepo.CloseOrder(tx, id)
	if err != nil {
		return false, err
	}
	return close, nil
}

func (s *WorkOrderServiceImpl) GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKey("all_request", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Error unmarshalling cached data"}
		}
		return result, 0, 0, nil
	}

	// Data not found in cache, proceed to database
	// Start database transaction
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Retrieve data from repository
	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllRequest(tx, filterCondition, pages)
	if repoErr != nil {

		errorResponse := &exceptions.BaseErrorResponse{Message: repoErr.Message}
		return nil, 0, 0, errorResponse
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Error marshalling data"}
	}

	if err := s.RedisClient.Set(context.Background(), cacheKey, jsonData, cacheExpiration).Err(); err != nil {
		fmt.Println("Error caching data:", err)
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKeyId("request_by_id", idwosn, idwos)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var request transactionworkshoppayloads.WorkOrderServiceRequest
		if err := json.Unmarshal([]byte(cachedData), &request); err != nil {
			return request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return request, nil
	} else if err != redis.Nil {

		return transactionworkshoppayloads.WorkOrderServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	request, repoErr := s.structWorkOrderRepo.GetRequestById(tx, idwosn, idwos)
	if repoErr != nil {
		return request, repoErr
	}

	cacheData, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		fmt.Println("Failed to marshal request data for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return request, nil
}

func (s *WorkOrderServiceImpl) UpdateRequest(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.UpdateRequest(tx, idwosn, idwos, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) AddRequest(id int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.AddRequest(tx, id, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) DeleteRequest(id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.DeleteRequest(tx, id, IdWorkorder)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKey("vehicle_service", filterCondition, pages)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &results); err != nil {
			return results, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return results, 0, 0, nil
	} else if err != redis.Nil {

		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllVehicleService(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	cacheData, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		fmt.Println("Failed to marshal results for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKeyId("vehicle_service", idwosn, idwos)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var result transactionworkshoppayloads.WorkOrderServiceVehicleRequest
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	} else if err != redis.Nil {

		return transactionworkshoppayloads.WorkOrderServiceVehicleRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, repoErr := s.structWorkOrderRepo.GetVehicleServiceById(tx, idwosn, idwos)
	if repoErr != nil {
		return result, repoErr
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) UpdateVehicleService(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.UpdateVehicleService(tx, idwosn, idwos, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) AddVehicleService(id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.AddVehicleService(tx, id, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) DeleteVehicleService(id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.DeleteVehicleService(tx, id, IdWorkorder)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) Submit(tx *gorm.DB, id int) (bool, string, *exceptions.BaseErrorResponse) {
	submit, newDocumentNumber, err := s.structWorkOrderRepo.Submit(tx, id)
	if err != nil {
		return false, "", err
	}
	return submit, newDocumentNumber, nil
}

func (s *WorkOrderServiceImpl) GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKey("detail_work_order", filterCondition, pages)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &results); err != nil {
			return results, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		totalPages := int(math.Ceil(float64(len(results)) / float64(pages.Limit)))
		totalRows := len(results)

		start := pages.Page * pages.Limit
		end := start + pages.Limit
		if end > totalRows {
			end = totalRows
		}
		results = results[start:end]

		return results, totalPages, totalRows, nil
	} else if err != redis.Nil {

		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllDetailWorkOrder(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	cacheData, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		fmt.Println("Failed to marshal results for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKeyId("work_order_detail_id", idwosn, idwos)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var result transactionworkshoppayloads.WorkOrderDetailRequest
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		return result, nil
	} else if err != redis.Nil {

		return transactionworkshoppayloads.WorkOrderDetailRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, repoErr := s.structWorkOrderRepo.GetDetailByIdWorkOrder(tx, idwosn, idwos)
	if repoErr != nil {
		return result, repoErr
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) UpdateDetailWorkOrder(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.UpdateDetailWorkOrder(tx, idwosn, idwos, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) AddDetailWorkOrder(id int, request transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.AddDetailWorkOrder(tx, id, request)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) DeleteDetailWorkOrder(id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.DeleteDetailWorkOrder(tx, id, IdWorkorder)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) NewBooking(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.NewBooking(tx, workOrderId, request)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	ctx := context.Background()
	cacheKey := generateCacheKey("all_booking", filterCondition, pages)

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {

		fmt.Println("Cache miss, querying database...")

		tx := s.DB.Begin()
		defer helper.CommitOrRollback(tx)

		results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllBooking(tx, filterCondition, pages)
		if repoErr != nil {
			return results, totalPages, totalRows, repoErr
		}

		cacheData, marshalErr := json.Marshal(results)
		if marshalErr == nil {
			s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
		} else {
			fmt.Println("Failed to marshal results for caching:", marshalErr)
		}

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

func (s *WorkOrderServiceImpl) GetBookingById(workOrderId int, id int) (transactionworkshoppayloads.WorkOrderBookingRequest, *exceptions.BaseErrorResponse) {

	cacheKey := generateCacheKeyId("booking", workOrderId, id)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {

		var result transactionworkshoppayloads.WorkOrderBookingRequest
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		return result, nil
	} else if err != redis.Nil {

		return transactionworkshoppayloads.WorkOrderBookingRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, repoErr := s.structWorkOrderRepo.GetBookingById(tx, workOrderId, id)
	if repoErr != nil {
		return result, repoErr
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, cacheExpiration)
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) SaveBooking(tx *gorm.DB, workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.SaveBooking(tx, workOrderId, id, request)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) SubmitBooking(tx *gorm.DB, workOrderId int, id int) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	submit, err := s.structWorkOrderRepo.SubmitBooking(tx, workOrderId, id)
	if err != nil {
		return false, err
	}
	return submit, nil
}

func (s *WorkOrderServiceImpl) VoidBooking(tx *gorm.DB, workOrderId int, id int) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	delete, err := s.structWorkOrderRepo.VoidBooking(tx, workOrderId, id)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) CloseBooking(tx *gorm.DB, workOrderId int, id int) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	close, err := s.structWorkOrderRepo.CloseBooking(tx, workOrderId, id)
	if err != nil {
		return false, err
	}
	return close, nil
}
