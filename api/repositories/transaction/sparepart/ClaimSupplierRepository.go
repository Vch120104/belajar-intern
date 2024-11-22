package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"gorm.io/gorm"
)

type ClaimSupplierRepository interface {
	InsertItemClaim(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse)
	InsertItemClaimDetail(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse)
	GetItemClaimById(db *gorm.DB, itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse)
}
