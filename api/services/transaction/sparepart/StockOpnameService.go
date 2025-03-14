package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type StockOpnameService interface {
	GetAllStockOpname([]utils.FilterCondition, pagination.Pagination, map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllStockOpnameDetail(pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetStockOpnameByStockOpnameSystemNumber(int) ([]transactionsparepartpayloads.GetStockOpnameByStockOpnameSystemNumberResponse, *exceptions.BaseErrorResponse)
	GetStockOpnameAllDetailByStockOpnameSystemNumber(int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitStockOpname(int, transactionsparepartpayloads.StockOpnameSubmitRequest) (bool, *exceptions.BaseErrorResponse)
	InsertStockOpname(transactionsparepartpayloads.StockOpnameInsertRequest) (bool, *exceptions.BaseErrorResponse)
	InsertStockOpnameDetail(transactionsparepartpayloads.StockOpnameInsertDetailRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateStockOpname(transactionsparepartpayloads.StockOpnameInsertRequest, int) (bool, *exceptions.BaseErrorResponse)
	UpdateStockOpnameDetail(transactionsparepartpayloads.StockOpnameUpdateDetailRequest, int) (bool, *exceptions.BaseErrorResponse)
	DeleteStockOpname(int) (bool, *exceptions.BaseErrorResponse)
	DeleteStockOpnameDetailByLineNumber(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteStockOpnameDetailBySystemNumber(int) (bool, *exceptions.BaseErrorResponse)
}
