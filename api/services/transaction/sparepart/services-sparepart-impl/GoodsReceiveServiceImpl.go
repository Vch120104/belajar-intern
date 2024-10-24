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

type GoodsReceiveServiceImpl struct {
	repository transactionsparepartrepository.GoodsReceiveRepository
	DB         *gorm.DB
	redis      *redis.Client
}

func NewGoodsReceiveServiceImpl(repository transactionsparepartrepository.GoodsReceiveRepository, db *gorm.DB, redis *redis.Client) transactionsparepartservice.GoodsReceiveService {
	return &GoodsReceiveServiceImpl{repository: repository, DB: db, redis: redis}

}
func (service *GoodsReceiveServiceImpl) GetAllGoodsReceive(filter []utils.FilterCondition, paginations pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.GetAllGoodsReceive(tx, filter, paginations)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) GetGoodsReceiveById(GoodsReceiveId int) (transactionsparepartpayloads.GoodsReceivesGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.GetGoodsReceiveById(tx, GoodsReceiveId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) InsertGoodsReceive(payloads transactionsparepartpayloads.GoodsReceiveInsertPayloads) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.InsertGoodsReceive(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) UpdateGoodsReceive(payloads transactionsparepartpayloads.GoodsReceiveUpdatePayloads, GoodsReceiveId int) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.UpdateGoodsReceive(tx, payloads, GoodsReceiveId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
