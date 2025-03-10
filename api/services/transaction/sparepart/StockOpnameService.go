package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type StockOpnameService interface {
	GetAllStockOpname([]utils.FilterCondition, pagination.Pagination, float64, map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLocationList([]utils.FilterCondition, pagination.Pagination, float64, string, string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetPersonInChargeList([]utils.FilterCondition, pagination.Pagination, float64) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemList(pagination.Pagination, string, string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOnGoingStockOpname(float64, float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse)
	InsertNewStockOpname(transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateOnGoingStockOpname(float64, transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse)
}
