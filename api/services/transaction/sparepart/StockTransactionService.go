package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
)

type StockTransactionService interface {
	StockTransactionInsert(payloads transactionsparepartpayloads.StockTransactionInsertPayloads) (bool, *exceptions.BaseErrorResponse)
}
