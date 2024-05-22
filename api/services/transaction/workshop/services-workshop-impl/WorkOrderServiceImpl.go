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

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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

func (s *WorkOrderServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) VehicleLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.VehicleLookup(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) New(tx *gorm.DB) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)

	results, err := s.structWorkOrderRepo.New(tx)
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, err
	}
	return results, nil
}

func (s *WorkOrderServiceImpl) NewStatus(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	statuses, err := s.structWorkOrderRepo.NewStatus(tx)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *WorkOrderServiceImpl) NewType(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	types, err := s.structWorkOrderRepo.NewType(tx)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (s *WorkOrderServiceImpl) GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WorkOrderServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	// Menggunakan "=" untuk menginisialisasi tx dengan transaksi yang dimulai
	defer helper.CommitOrRollback(tx)

	// Panggil metode Save dengan menyediakan transaksi dan permintaan WorkOrder
	save, err := s.structWorkOrderRepo.Save(tx, request)
	if err != nil {
		return false, err
	}

	// Mengembalikan hasil penyimpanan dan nilai nil untuk ErrorResponse
	return save, nil
}

func (s *WorkOrderServiceImpl) Submit(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.Submit(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) Void(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.Void(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.CloseOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}
