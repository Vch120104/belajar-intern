package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"gorm.io/gorm"
)

type BinningListServiceImpl struct {
	DB         *gorm.DB
	repository transactionsparepartrepository.BinningListRepository
}

func NewBinningListServiceImpl(db *gorm.DB, repository transactionsparepartrepository.BinningListRepository) transactionsparepartservice.BinningListService {
	return &BinningListServiceImpl{
		DB:         db,
		repository: repository,
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
