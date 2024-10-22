package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	utils "after-sales/api/utils"
	"log"
	"net/http"

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

// Function to generate document number
func (s *WorkOrderServiceImpl) GenerateDocumentNumber(workOrderId int) (string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	documentNumber, err := s.structWorkOrderRepo.GenerateDocumentNumber(tx, workOrderId)
	if err != nil {
		return "", err
	}
	log.Printf("Document number from repository: %s", documentNumber)
	return documentNumber, nil
}

func (s *WorkOrderServiceImpl) NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	statuses, err := s.structWorkOrderRepo.NewStatus(tx, filter)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *WorkOrderServiceImpl) AddStatus(request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddStatus(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateStatus(id int, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateStatus(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil
}

func (s *WorkOrderServiceImpl) DeleteStatus(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteStatus(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewType(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	types, err := s.structWorkOrderRepo.NewType(tx, filter)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (s *WorkOrderServiceImpl) AddType(request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddType(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateType(id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateType(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil
}

func (s *WorkOrderServiceImpl) DeleteType(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteType(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewBill(filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	bills, err := s.structWorkOrderRepo.NewBill(tx, filter)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) AddBill(request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddBill(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateBill(id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateBill(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil

}

func (s *WorkOrderServiceImpl) DeleteBill(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteBill(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewDropPoint() ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	dropPoints, err := s.structWorkOrderRepo.NewDropPoint(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return dropPoints, nil
}

func (s *WorkOrderServiceImpl) NewVehicleBrand() ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	brands, err := s.structWorkOrderRepo.NewVehicleBrand(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (s *WorkOrderServiceImpl) NewVehicleModel(brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	models, err := s.structWorkOrderRepo.NewVehicleModel(tx, brandId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *WorkOrderServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) New(request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.New(tx, request)
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.structWorkOrderRepo.GetById(tx, id, pages)
	if repoErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, repoErr
	}

	return results, nil
}

func (s *WorkOrderServiceImpl) Save(request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	save, err := s.structWorkOrderRepo.Save(tx, request, workOrderId)
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) Void(workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	delete, err := s.structWorkOrderRepo.Void(tx, workOrderId)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) CloseOrder(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	close, err := s.structWorkOrderRepo.CloseOrder(tx, id)
	if err != nil {
		return false, err
	}
	return close, nil
}

func (s *WorkOrderServiceImpl) GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	// Data not found in cache, proceed to database
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllRequest(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	request, repoErr := s.structWorkOrderRepo.GetRequestById(tx, idwosn, idwos)

	if repoErr != nil {
		return request, repoErr
	}

	return request, nil
}

func (s *WorkOrderServiceImpl) UpdateRequest(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	update, err := s.structWorkOrderRepo.UpdateRequest(tx, idwosn, idwos, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderRequestDescription{}, err
	}

	return update, nil
}

func (s *WorkOrderServiceImpl) AddRequest(id int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.AddRequest(tx, id, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderRequestDescription{}, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) DeleteRequest(id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	delete, err := s.structWorkOrderRepo.DeleteRequest(tx, id, IdWorkorder)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllVehicleService(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.structWorkOrderRepo.GetVehicleServiceById(tx, idwosn, idwos)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) UpdateVehicleService(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	update, err := s.structWorkOrderRepo.UpdateVehicleService(tx, idwosn, idwos, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionworkshopentities.WorkOrderServiceVehicle{}, err
	}

	return update, nil
}

func (s *WorkOrderServiceImpl) AddVehicleService(id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.AddVehicleService(tx, id, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderServiceVehicle{}, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) DeleteVehicleService(id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	delete, err := s.structWorkOrderRepo.DeleteVehicleService(tx, id, IdWorkorder)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}

	return delete, nil
}

func (s *WorkOrderServiceImpl) Submit(id int) (bool, string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	submit, newDocumentNumber, err := s.structWorkOrderRepo.Submit(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, "", err
	}

	return submit, newDocumentNumber, nil
}

func (s *WorkOrderServiceImpl) GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllDetailWorkOrder(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil

}

func (s *WorkOrderServiceImpl) GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.structWorkOrderRepo.GetDetailByIdWorkOrder(tx, idwosn, idwos)

	if repoErr != nil {
		if repoErr.StatusCode == http.StatusNotFound {
			return result, repoErr
		}
		return result, repoErr
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) UpdateDetailWorkOrder(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	update, err := s.structWorkOrderRepo.UpdateDetailWorkOrder(tx, idwosn, idwos, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}

	return update, nil
}

func (s *WorkOrderServiceImpl) AddDetailWorkOrder(id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	submit, err := s.structWorkOrderRepo.AddDetailWorkOrder(tx, id, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}

	return submit, nil
}

func (s *WorkOrderServiceImpl) DeleteDetailWorkOrder(id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	delete, err := s.structWorkOrderRepo.DeleteDetailWorkOrder(tx, id, IdWorkorder)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}

	return delete, nil
}

func (s *WorkOrderServiceImpl) GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllBooking(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) GetBookingById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.structWorkOrderRepo.GetBookingById(tx, workOrderId, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) NewBooking(request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.NewBooking(tx, request)
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) SaveBooking(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.SaveBooking(tx, workOrderId, id, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) GetAllAffiliated(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderRepo.GetAllAffiliated(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil

}

func (s *WorkOrderServiceImpl) GetAffiliatedById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.structWorkOrderRepo.GetAffiliatedById(tx, workOrderId, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *WorkOrderServiceImpl) NewAffiliated(workOrderId int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.NewAffiliated(tx, workOrderId, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) SaveAffiliated(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.SaveAffiliated(tx, workOrderId, id, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) DeleteRequestMultiId(workOrderId int, id []int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	deletemultiid, err := s.structWorkOrderRepo.DeleteRequestMultiId(tx, workOrderId, id)
	if err != nil {
		return false, err
	}

	return deletemultiid, nil
}

func (s *WorkOrderServiceImpl) DeleteVehicleServiceMultiId(workOrderId int, id []int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	deletemultiid, err := s.structWorkOrderRepo.DeleteVehicleServiceMultiId(tx, workOrderId, id)
	if err != nil {
		return false, err
	}

	return deletemultiid, nil
}

func (s *WorkOrderServiceImpl) DeleteDetailWorkOrderMultiId(workOrderId int, id []int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	deletemultiid, err := s.structWorkOrderRepo.DeleteDetailWorkOrderMultiId(tx, workOrderId, id)
	if err != nil {
		return false, err
	}

	return deletemultiid, nil
}

func (s *WorkOrderServiceImpl) ChangeBillTo(workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	change, err := s.structWorkOrderRepo.ChangeBillTo(tx, workOrderId, request)
	if err != nil {
		return false, err
	}

	return change, nil
}

func (s *WorkOrderServiceImpl) ChangePhoneNo(workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (*transactionworkshoppayloads.ChangePhoneNoRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	updatedPayload, err := s.structWorkOrderRepo.ChangePhoneNo(tx, workOrderId, request)
	if err != nil {
		return nil, err
	}

	return updatedPayload, nil
}

func (s *WorkOrderServiceImpl) ConfirmPrice(workOrderId int, idwos []int, request transactionworkshoppayloads.WorkOrderConfirmPriceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	confirm, err := s.structWorkOrderRepo.ConfirmPrice(tx, workOrderId, idwos, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}

	return confirm, nil
}

func (s *WorkOrderServiceImpl) NewTrxType(filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	bills, err := s.structWorkOrderRepo.NewTrxType(tx, filter)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) AddTrxType(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddTrxType(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateTrxType(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateTrxType(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil

}

func (s *WorkOrderServiceImpl) DeleteTrxType(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteTrxType(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewLineType(filter []utils.FilterCondition) ([]transactionworkshoppayloads.Linetype, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	bills, err := s.structWorkOrderRepo.NewLineType(tx, filter)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) AddLineType(request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddLineType(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateLineType(id int, request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateLineType(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil

}

func (s *WorkOrderServiceImpl) DeleteLineType(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteLineType(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewTrxTypeSo(filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	bills, err := s.structWorkOrderRepo.NewTrxTypeSo(tx, filter)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) AddTrxTypeSo(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddTrxTypeSo(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateTrxTypeSo(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateTrxTypeSo(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil

}

func (s *WorkOrderServiceImpl) DeleteTrxTypeSo(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteTrxTypeSo(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) NewJobType(filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderJobType, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	bills, err := s.structWorkOrderRepo.NewJobType(tx, filter)
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (s *WorkOrderServiceImpl) AddJobType(request transactionworkshoppayloads.WorkOrderJobType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	save, err := s.structWorkOrderRepo.AddJobType(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) UpdateJobType(id int, request transactionworkshoppayloads.WorkOrderJobType) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	update, err := s.structWorkOrderRepo.UpdateJobType(tx, id, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return update, nil

}

func (s *WorkOrderServiceImpl) DeleteJobType(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	delete, err := s.structWorkOrderRepo.DeleteJobType(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) DeleteCampaign(workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	delete, err := s.structWorkOrderRepo.DeleteCampaign(tx, workOrderId)
	if err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, err
	}
	return delete, nil
}

func (s *WorkOrderServiceImpl) AddContractService(workOrderId int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.AddContractService(tx, workOrderId, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) AddGeneralRepairPackage(workOrderId int, request transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.AddGeneralRepairPackage(tx, workOrderId, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) AddFieldAction(workOrderId int, request transactionworkshoppayloads.WorkOrderFieldActionRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	save, err := s.structWorkOrderRepo.AddFieldAction(tx, workOrderId, request)
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, err
	}
	return save, nil
}
