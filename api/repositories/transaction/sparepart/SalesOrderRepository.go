package transactionsparepartrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

// SalesOrderRepository
type SalesOrderRepository interface {
	GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptionsss_test.BaseErrorResponse)
}
