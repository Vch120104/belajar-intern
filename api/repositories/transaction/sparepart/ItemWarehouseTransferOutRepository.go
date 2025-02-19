package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemWarehouseTransferOutRepository interface {
	InsertHeader(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse)
	InsertDetail(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	InsertDetailFromReceipt(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	GetTransferOutById(*gorm.DB, int) (transactionsparepartpayloads.GetTransferOutByIdResponse, *exceptions.BaseErrorResponse)
	GetAllTransferOut(*gorm.DB, []utils.FilterCondition, map[string]string, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllTransferOutDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitTransferOut(*gorm.DB, int) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse)
	UpdateTransferOutDetail(*gorm.DB, transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest, int) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse)
	DeleteTransferOutDetail(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeleteTransferOut(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
