package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type ItemWarehouseTransferOutService interface {
	InsertHeader(transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse)
	InsertDetail(transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	InsertDetailFromReceipt(transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	GetTransferOutById(int) (transactionsparepartpayloads.GetTransferOutByIdResponse, *exceptions.BaseErrorResponse)
	GetAllTransferOut([]utils.FilterCondition, map[string]string, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllTransferOutDetail(int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitTransferOut(int) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse)
	UpdateTransferOutDetail(transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest, int) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	DeleteTransferOutDetail(int) (bool, *exceptions.BaseErrorResponse)
	DeleteTransferOut(int) (bool, *exceptions.BaseErrorResponse)
}
