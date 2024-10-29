package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"gorm.io/gorm"
)

type StockTransactionReasonRepository interface {
	GetStockTransactionReasonByCode(db *gorm.DB, Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
}
