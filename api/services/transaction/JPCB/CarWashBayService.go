package transactionjpcbservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type BayMasterService interface {
	GetAllCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllActiveCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDeactiveCarWashBay(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	ChangeStatusCarWashBay(request transactionjpcbpayloads.BayMasterUpdateRequest) (bool, *exceptions.BaseErrorResponse)
}
