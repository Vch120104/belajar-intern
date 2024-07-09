package transactionsparepartservice

import (
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type SalesOrderService interface {
	GetSalesOrderByID(tx *gorm.DB, id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse)
}
