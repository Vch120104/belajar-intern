package masterservice

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type StockTransactionTypeService interface {
	GetStockTransactionTypeByCode(Code string) (masterentities.StockTransactionType, *exceptions.BaseErrorResponse)
	GetAllStockTransactionType([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
