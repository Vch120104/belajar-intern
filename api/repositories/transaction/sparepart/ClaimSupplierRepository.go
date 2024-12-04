package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type ClaimSupplierRepository interface {
	InsertItemClaim(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse)
	InsertItemClaimDetail(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse)
	GetItemClaimById(db *gorm.DB, itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse)
	GetItemClaimDetailByHeaderId(db *gorm.DB, Paginations pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitItemClaim(db *gorm.DB, claimId int) (bool, *exceptions.BaseErrorResponse)
	GetAllItemClaim(db *gorm.DB, page pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
