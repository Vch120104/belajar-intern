package transactionworkshopserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

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

func (s *BookingEstimationServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structBookingEstimationRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BookingEstimationServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	results, err := s.structBookingEstimationRepo.New(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) GetById(id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structBookingEstimationRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	// Menggunakan "=" untuk menginisialisasi tx dengan transaksi yang dimulai
	defer helper.CommitOrRollback(tx)

	// Panggil metode Save dengan menyediakan transaksi dan permintaan WorkOrder
	save, err := s.structBookingEstimationRepo.Save(tx, request)
	if err != nil {
		return false, err
	}

	// Mengembalikan hasil penyimpanan dan nilai nil untuk ErrorResponse
	return save, nil
}

func (s *BookingEstimationServiceImpl) Submit(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.Submit(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) Void(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.Void(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structBookingEstimationRepo.CloseOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}
