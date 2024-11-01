package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"gorm.io/gorm"
)

type StockTransactionRepository interface {
	StockTransactionInsert(db *gorm.DB, payloads transactionsparepartpayloads.StockTransactionInsertPayloads) (bool, *exceptions.BaseErrorResponse)
}
