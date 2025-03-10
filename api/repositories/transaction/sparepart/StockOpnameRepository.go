package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type StockOpnameRepository interface {
	GetAllStockOpname(*gorm.DB, []utils.FilterCondition, pagination.Pagination, float64, map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLocationList(*gorm.DB, []utils.FilterCondition, pagination.Pagination, float64, string, string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetPersonInChargeList(*gorm.DB, []utils.FilterCondition, pagination.Pagination, float64) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemList(*gorm.DB, pagination.Pagination, string, string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOnGoingStockOpname(*gorm.DB, float64, float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse)
	InsertNewStockOpname(*gorm.DB, transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateOnGoingStockOpname(*gorm.DB, float64, transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse)
}
