package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
)

type BinningListService interface {
	GetBinningListById(BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse)
}
