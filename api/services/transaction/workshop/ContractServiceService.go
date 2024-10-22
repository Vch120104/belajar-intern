package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ContractServiceService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
