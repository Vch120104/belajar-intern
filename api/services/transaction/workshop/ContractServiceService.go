package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ContractServiceService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse)
}
