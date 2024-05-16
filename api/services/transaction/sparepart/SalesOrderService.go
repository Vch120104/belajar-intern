package transactionsparepartservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type SalesOrderService interface {
	GetSalesOrderByID(tx *gorm.DB, id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptionsss_test.BaseErrorResponse)
}
