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
func (service *BinningListServiceImpl) InsertBinningListHeader(payloads transactionsparepartpayloads.BinningListInsertPayloads) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.InsertBinningListHeader(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *BinningListServiceImpl) UpdateBinningListHeader(payloads transactionsparepartpayloads.BinningListSavePayload) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.UpdateBinningListHeader(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) GetBinningListDetailById(BinningDetailSystemNumber int) (transactionsparepartpayloads.BinningListGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.GetBinningListDetailById(tx, BinningDetailSystemNumber)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) GetAllBinningListDetailWithPagination(filter []utils.FilterCondition, pagination pagination.Pagination, binningListId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.GetAllBinningListDetailWithPagination(tx, filter, pagination, binningListId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *BinningListServiceImpl) InsertBinningListDetail(payloads transactionsparepartpayloads.BinningListDetailPayloads) (transactionsparepartentities.BinningStockDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.InsertBinningListDetail(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) UpdateBinningListDetail(payloads transactionsparepartpayloads.BinningListDetailUpdatePayloads) (transactionsparepartentities.BinningStockDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.UpdateBinningListDetail(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *BinningListServiceImpl) SubmitBinningList(BinningId int) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.SubmitBinningList(tx, BinningId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) DeleteBinningList(BinningId int) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.DeleteBinningList(tx, BinningId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *BinningListServiceImpl) DeleteBinningListDetailMultiId(binningDetailMultiId string) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.repository.DeleteBinningListDetailMultiId(tx, binningDetailMultiId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	//tx.Rollback()
	return result, nil
}
