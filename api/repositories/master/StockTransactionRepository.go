package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type StockTransactionTypeRepository interface {
	GetStockTransactionTypeByCode(db *gorm.DB, Code string) (masterentities.StockTransactionType, *exceptions.BaseErrorResponse)
	GetAllStockTransactionType(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
