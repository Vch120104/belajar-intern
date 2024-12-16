package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLocationTransferRepository interface {
	GetAllItemLocationTransfer(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLocationTransferById(tx *gorm.DB, id int) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	InsertItemLocationTransfer(tx *gorm.DB, request transactionsparepartpayloads.InsertItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	UpdateItemLocationTransfer(tx *gorm.DB, id int, request transactionsparepartpayloads.UpdateItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	AcceptItemLocationTransfer(tx *gorm.DB, id int, request transactionsparepartpayloads.AcceptItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	RejectItemLocationTransfer(tx *gorm.DB, id int, request transactionsparepartpayloads.RejectItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	DeleteItemLocationTransfer(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)

	GetAllItemLocationTransferDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	InsertItemLocationTransferDetail(tx *gorm.DB, request transactionsparepartpayloads.InsertItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse)
	UpdateItemLocationTransferDetail(tx *gorm.DB, id int, request transactionsparepartpayloads.UpdateItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse)
	DeleteItemLocationTransferDetail(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse)
}
