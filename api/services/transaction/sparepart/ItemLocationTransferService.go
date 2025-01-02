package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type ItemLocationTransferService interface {
	GetAllItemLocationTransfer(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLocationTransferById(id int) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	InsertItemLocationTransfer(request transactionsparepartpayloads.InsertItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	UpdateItemLocationTransfer(id int, request transactionsparepartpayloads.UpdateItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	AcceptItemLocationTransfer(id int, request transactionsparepartpayloads.AcceptItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	RejectItemLocationTransfer(id int, request transactionsparepartpayloads.RejectItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	SubmitItemLocationTransfer(id int, request transactionsparepartpayloads.SubmitItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse)
	DeleteItemLocationTransfer(id int) (bool, *exceptions.BaseErrorResponse)

	GetAllItemLocationTransferDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLocationTransferDetailById(id int) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse)
	InsertItemLocationTransferDetail(request transactionsparepartpayloads.InsertItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse)
	UpdateItemLocationTransferDetail(id int, request transactionsparepartpayloads.UpdateItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse)
	DeleteItemLocationTransferDetail(ids []int) (bool, *exceptions.BaseErrorResponse)
}
