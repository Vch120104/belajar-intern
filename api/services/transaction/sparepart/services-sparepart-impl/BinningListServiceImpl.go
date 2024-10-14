package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BinningListServiceImpl struct {
	DB         *gorm.DB
	repository transactionsparepartrepository.BinningListRepository
	redis      *redis.Client
}

func NewBinningListServiceImpl(repository transactionsparepartrepository.BinningListRepository, db *gorm.DB, redisclient *redis.Client) transactionsparepartservice.BinningListService {
	return &BinningListServiceImpl{
		DB:         db,
		repository: repository,
		redis:      redisclient,
	}
}

func (service *BinningListServiceImpl) GetBinningListById(BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	//rdbb := service.redis
	//errs := rdbb.Set(nil, "dasa", "123", 0)
	//data := rdbb.Get("dasa")
	result, err := service.repository.GetBinningListById(tx, BinningStockId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) GetAllBinningListWithPagination(filter []utils.FilterCondition, pagination pagination.Pagination, ctx context.Context) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	rdb := service.redis
	result, err := service.repository.GetAllBinningListWithPagination(tx, rdb, filter, pagination, ctx)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
