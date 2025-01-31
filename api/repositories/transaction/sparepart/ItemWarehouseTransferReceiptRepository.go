package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type ItemWarehouseTransferReceiptRepository interface {
	Accept(*gorm.DB, int, transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	Reject(*gorm.DB, int, transactionsparepartpayloads.RejectWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
}
