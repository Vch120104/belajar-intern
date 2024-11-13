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
func (service *GoodsReceiveServiceImpl) InsertGoodsReceiveDetail(payloads transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads) (transactionsparepartentities.GoodsReceiveDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.InsertGoodsReceiveDetail(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) UpdateGoodsReceiveDetail(payloads transactionsparepartpayloads.GoodsReceiveDetailUpdatePayloads, DetailId int) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.UpdateGoodsReceiveDetail(tx, payloads, DetailId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) LocationItemGoodsReceive(filter []utils.FilterCondition, PaginationParams pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.LocationItemGoodsReceive(tx, filter, PaginationParams)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *GoodsReceiveServiceImpl) SubmitGoodsReceive(GoodsReceiveId int) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.SubmitGoodsReceive(tx, GoodsReceiveId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *GoodsReceiveServiceImpl) DeleteGoodsReceive(goodsReceivesId int) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.DeleteGoodsReceive(tx, goodsReceivesId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
