package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ClaimSupplierServiceImpl struct {
	claimRepository transactionsparepartrepository.ClaimSupplierRepository
	DB              *gorm.DB
	rdb             *redis.Client
}

func NewClaimSupplierServiceImpl(repo transactionsparepartrepository.ClaimSupplierRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.ClaimSupplierService {
	return &ClaimSupplierServiceImpl{
		claimRepository: repo,
		DB:              db,
		rdb:             rdb,
	}
}
func (service *ClaimSupplierServiceImpl) InsertItemClaim(payload transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.claimRepository.InsertItemClaim(tx, payload)
	//tx.Rollback()
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *ClaimSupplierServiceImpl) InsertItemClaimDetail(payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.claimRepository.InsertItemClaimDetail(tx, payloads)
	//tx.Rollback()
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *ClaimSupplierServiceImpl) GetItemClaimById(itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.claimRepository.GetItemClaimById(tx, itemClaimId)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
