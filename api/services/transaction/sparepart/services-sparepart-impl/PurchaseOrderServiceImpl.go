package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
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

func NewPurchaseOrderService(PuchaseOrderRepo transactionsparepartrepository.PurchaseOrderRepository, db *gorm.DB, rediscliend *redis.Client) transactionsparepartservice.PurchaseOrderService {
	return &PurchaseOrderServiceImpl{
		PurchaseOrderRepo: PuchaseOrderRepo,
		DB:                db,
		redis:             rediscliend,
	}
}

func (service *PurchaseOrderServiceImpl) GetAllPurchaseOrder(filter []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := service.DB.Begin()
	result, err := service.PurchaseOrderRepo.GetAllPurchaseOrder(tx, filter, page)
	if err != nil {
		return result, err
	}
	return result, nil
}
