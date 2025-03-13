package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type StockOpnameRepository interface {
	GetAllStockOpname(*gorm.DB, []utils.FilterCondition, pagination.Pagination, map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllStockOpnameDetail(*gorm.DB, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetStockOpnameByStockOpnameSystemNumber(*gorm.DB, int) (
		[]transactionsparepartpayloads.GetStockOpnameByStockOpnameSystemNumberResponse, *exceptions.BaseErrorResponse)
	GetStockOpnameAllDetailByStockOpnameSystemNumber(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	InsertStockOpname(*gorm.DB, transactionsparepartpayloads.StockOpnameInsertRequest) (bool, *exceptions.BaseErrorResponse)
	SubmitStockOpname(*gorm.DB, int, transactionsparepartpayloads.StockOpnameSubmitRequest) (bool, *exceptions.BaseErrorResponse)
	InsertStockOpnameDetail(*gorm.DB, transactionsparepartpayloads.StockOpnameInsertDetailRequest, int) (bool, *exceptions.BaseErrorResponse)
	UpdateStockOpname(*gorm.DB, transactionsparepartpayloads.StockOpnameInsertRequest, int) (bool, *exceptions.BaseErrorResponse)
	UpdateStockOpnameDetail(*gorm.DB, transactionsparepartpayloads.StockOpnameUpdateDetailRequest, int) (bool, *exceptions.BaseErrorResponse)
	DeleteStockOpname(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
