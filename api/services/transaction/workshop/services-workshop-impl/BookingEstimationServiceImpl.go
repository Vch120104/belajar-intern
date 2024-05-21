package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookingEstimationServiceImpl struct {
	structBookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository
	DB                          *gorm.DB
	RedisClient                 *redis.Client // Redis client
}

func OpenBookingEstimationServiceImpl(bookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.BookingEstimationService {
	return &BookingEstimationServiceImpl{
		structBookingEstimationRepo: bookingEstimationRepo,
		DB:                          db,
		RedisClient:                 redisClient,
	}
}

func (s *BookingEstimationServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structBookingEstimationRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BookingEstimationServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	_, err := s.structBookingEstimationRepo.New(tx, request)
	if err != nil {
		return exceptions.NewBaseErrorResponse(http.StatusInternalServerError, "Failed to create new record", err)
	}
	return nil
}

func (s *BookingEstimationServiceImpl) GetById(id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structBookingEstimationRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	_, err := s.structBookingEstimationRepo.Save(tx, request)
	if err != nil {
		return exceptions.NewBaseErrorResponse(http.StatusInternalServerError, "Failed to save record", err)
	}
	return nil
}

func (s *BookingEstimationServiceImpl) Submit(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.Submit(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) Void(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.Void(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.CloseOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}
