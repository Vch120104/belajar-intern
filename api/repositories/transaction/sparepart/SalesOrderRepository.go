package transactionsparepartrepository

import (
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

// SalesOrderRepository
type SalesOrderRepository interface {
	GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse)
}
