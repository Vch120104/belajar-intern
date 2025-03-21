package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ContractServiceService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse)
	Save(payload transactionworkshoppayloads.ContractServiceInsert) (transactionworkshoppayloads.ContractServiceInsert, *exceptions.BaseErrorResponse)
	Void(Id int) (bool, *exceptions.BaseErrorResponse)
	Submit(Id int) (bool, *exceptions.BaseErrorResponse)
}
