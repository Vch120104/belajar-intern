package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PurchaseOrderServiceImpl struct {
	PurchaseOrderRepo transactionsparepartrepository.PurchaseOrderRepository
	DB                *gorm.DB
	redis             *redis.Client
}

func NewPurchaseOrderService(PuchaseOrderRepo transactionsparepartrepository.PurchaseOrderRepository, db *gorm.DB, redisclient *redis.Client) transactionsparepartservice.PurchaseOrderService {
	return &PurchaseOrderServiceImpl{
		PurchaseOrderRepo: PuchaseOrderRepo,
		DB:                db,
		redis:             redisclient,
	}
}

func (service *PurchaseOrderServiceImpl) GetAllPurchaseOrder(filter []utils.FilterCondition, page pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.GetAllPurchaseOrder(tx, filter, page, DateParams)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *PurchaseOrderServiceImpl) GetByIdPurchaseOrder(i int) (transactionsparepartpayloads.PurchaseOrderGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.GetByIdPurchaseOrder(tx, i)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *PurchaseOrderServiceImpl) GetByIdPurchaseOrderDetail(id int, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.GetAllDetailByHeaderId(tx, id, page)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *PurchaseOrderServiceImpl) NewPurchaseOrderHeader(responses transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.NewPurchaseOrderHeader(tx, responses)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *PurchaseOrderServiceImpl) UpdatePurchaseOrderHeader(id int, responses transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.UpdatePurchaseOrderHeader(tx, id, responses)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *PurchaseOrderServiceImpl) GetPurchaseOrderDetailById(id int) (transactionsparepartpayloads.PurchaseOrderGetDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.GetPurchaseOrderDetailById(tx, id)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *PurchaseOrderServiceImpl) NewPurchaseOrderDetail(payloads transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.NewPurchaseOrderDetail(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *PurchaseOrderServiceImpl) DeletePurchaseOrderDetailMultiId(multiid string) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.DeletePurchaseOrderDetailMultiId(tx, multiid)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *PurchaseOrderServiceImpl) SavePurchaseOrderDetail(payloads transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.SavePurchaseOrderDetail(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
