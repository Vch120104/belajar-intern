package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"gorm.io/gorm"
)

type SalesOrderRepository interface {
	InsertSalesOrderHeader(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderInsertHeaderPayload) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse)
	GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse)
}
