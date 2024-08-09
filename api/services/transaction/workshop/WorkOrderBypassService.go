package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type WorkOrderBypassService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse)
	Bypass(request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshopentities.WorkOrderQualityControl, *exceptions.BaseErrorResponse)
}
