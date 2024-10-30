package masterservice

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type StockTransactionReasonService interface {
	GetStockTransactionReasonByCode(Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	InsertStockTransactionReason(payloads masterpayloads.StockTransactionReasonInsertPayloads) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	GetStockTransactionReasonById(id int) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse)
	GetAllStockTransactionReason([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
