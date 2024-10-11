package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"gorm.io/gorm"
)

type BinningListRepository interface {
	GetBinningListById(db *gorm.DB, BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse)
}
