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
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceRequestServiceImpl struct {
	ServiceRequestRepository transactionworkshoprepository.ServiceRequestRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenServiceRequestServiceImpl(ServiceRequestRepo transactionworkshoprepository.ServiceRequestRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceRequestService {
	return &ServiceRequestServiceImpl{
		ServiceRequestRepository: ServiceRequestRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

// Function to generate document service number
func (s *ServiceRequestServiceImpl) GenerateDocumentNumberServiceRequest(ServiceRequestId int) (string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	documentNumber, err := s.ServiceRequestRepository.GenerateDocumentNumberServiceRequest(tx, ServiceRequestId)
	if err != nil {
		return "", err
	}
	log.Printf("Document number from repository: %s", documentNumber)
	return documentNumber, nil
}

func (s *ServiceRequestServiceImpl) NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	statuses, err := s.ServiceRequestRepository.NewStatus(tx, filter)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *ServiceRequestServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.ServiceRequestRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *ServiceRequestServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {
	cacheKey := utils.GenerateCacheKeyIds("service_request_id", id)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result transactionworkshoppayloads.ServiceRequestResponse
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	} else if err != redis.Nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.ServiceRequestRepository.GetById(tx, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		if err := s.RedisClient.Set(ctx, cacheKey, cacheData, utils.CacheExpiration).Err(); err != nil {
			fmt.Println("Failed to set cache:", err)
		}
	}

	return result, nil
}

func (s *ServiceRequestServiceImpl) New(request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {

	ctx := context.Background()
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	save, err := s.ServiceRequestRepository.New(tx, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, err
	}

	utils.RefreshCaches(ctx, "service_request")

	return save, nil
}

func (s *ServiceRequestServiceImpl) Save(id int, request transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	save, err := s.ServiceRequestRepository.Save(tx, id, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, err
	}

	cacheKey := utils.GenerateCacheKeyIds("service_request_id", id)
	if err := s.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Println("Failed to delete cache:", err)
	}

	cacheData, marshalErr := json.Marshal(save)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		if err := s.RedisClient.Set(ctx, cacheKey, cacheData, utils.CacheExpiration).Err(); err != nil {
			fmt.Println("Failed to set cache:", err)
		}
	}

	return save, nil
}

func (s *ServiceRequestServiceImpl) Submit(id int) (bool, string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	submit, newDocumentNumber, err := s.ServiceRequestRepository.Submit(tx, id)
	if err != nil {
		return false, "", err
	}

	return submit, newDocumentNumber, nil
}

func (s *ServiceRequestServiceImpl) Void(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	void, err := s.ServiceRequestRepository.Void(tx, id)
	if err != nil {
		return false, err
	}
	return void, nil
}

func (s *ServiceRequestServiceImpl) CloseOrder(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	closeOrder, err := s.ServiceRequestRepository.CloseOrder(tx, id)
	if err != nil {
		return false, err
	}
	return closeOrder, nil
}

func (s *ServiceRequestServiceImpl) GetAllServiceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.ServiceRequestRepository.GetAllServiceDetail(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil

}

func (s *ServiceRequestServiceImpl) GetServiceDetailById(idsys int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.ServiceRequestRepository.GetServiceDetailById(tx, idsys)
	if repoErr != nil {
		if repoErr.StatusCode == http.StatusNotFound {
			return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshoppayloads.ServiceDetailResponse{}, repoErr
	}

	return results, nil

}

func (s *ServiceRequestServiceImpl) AddServiceDetail(idsys int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	success, err := s.ServiceRequestRepository.AddServiceDetail(tx, idsys, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, err
	}

	return success, nil
}

func (s *ServiceRequestServiceImpl) UpdateServiceDetail(idsys int, idservice int, request transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	save, err := s.ServiceRequestRepository.UpdateServiceDetail(tx, idsys, idservice, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, err
	}

	return save, nil
}

func (s *ServiceRequestServiceImpl) DeleteServiceDetail(idsys int, idservice int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	delete, err := s.ServiceRequestRepository.DeleteServiceDetail(tx, idsys, idservice)
	if err != nil {
		return false, err
	}

	return delete, nil
}

func (s *ServiceRequestServiceImpl) DeleteServiceDetailMultiId(idsys int, idservice []int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	deletemultiid, err := s.ServiceRequestRepository.DeleteServiceDetailMultiId(tx, idsys, idservice)
	if err != nil {
		return false, err
	}

	return deletemultiid, nil
}
