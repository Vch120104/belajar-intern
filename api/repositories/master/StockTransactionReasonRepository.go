package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type StockTransactionReasonRepository interface {
	GetStockTransactionReasonByCode(db *gorm.DB, Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	InsertStockTransactionReason(db *gorm.DB, payloads masterpayloads.StockTransactionReasonInsertPayloads) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	GetStockTransactionReasonById(db *gorm.DB, id int) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	GetAllStockTransactionReason(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
