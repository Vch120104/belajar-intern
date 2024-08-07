package transactionjpcbservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type BayMasterService interface {
	GetAllBayMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllActiveBayCarWashScreen(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDeactiveBayCarWashScreen(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
}
