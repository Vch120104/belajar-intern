package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type ClaimSupplierService interface {
	InsertItemClaim(payload transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse)
	InsertItemClaimDetail(payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse)
	GetItemClaimById(itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse)
	GetItemClaimDetailByHeaderId(Paginations pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitItemClaim(claimId int) (bool, *exceptions.BaseErrorResponse)
	GetAllItemClaim(page pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
