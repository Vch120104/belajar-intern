package masterservice

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
)

type StockTransactionReasonService interface {
	GetStockTransactionReasonByCode(Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
}
