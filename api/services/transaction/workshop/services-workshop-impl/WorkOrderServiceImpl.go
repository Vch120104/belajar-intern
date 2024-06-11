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

func (s *WorkOrderServiceImpl) CampaignLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.CampaignLookup(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.New(tx, request)
	if err != nil {
		return false, err
	}
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
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WorkOrderServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	// Start a new transaction
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.Save(tx, request, workOrderId)
	if err != nil {
		return false, err
	}
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
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAllRequest(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetRequestById(tx, idwosn, idwos)
	if err != nil {
		return results, err
	}
	return results, nil
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
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAllVehicleService(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetVehicleServiceById(tx, idwosn, idwos)
	if err != nil {
		return results, err
	}
	return results, nil
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
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAllDetailWorkOrder(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetDetailByIdWorkOrder(tx, idwosn, idwos)
	if err != nil {
		return results, err
	}
	return results, nil
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
