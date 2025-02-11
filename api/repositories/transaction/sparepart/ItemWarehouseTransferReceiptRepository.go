package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemWarehouseTransferReceiptRepository interface {
	GetAll(*gorm.DB, pagination.Pagination, []utils.FilterCondition, map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	Accept(*gorm.DB, int, transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	Reject(*gorm.DB, int, transactionsparepartpayloads.RejectWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
}
