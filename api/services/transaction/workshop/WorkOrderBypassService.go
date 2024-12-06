package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type WorkOrderBypassService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse)
	Bypass(id int, request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshoppayloads.WorkOrderBypassResponseDetail, *exceptions.BaseErrorResponse)
}
