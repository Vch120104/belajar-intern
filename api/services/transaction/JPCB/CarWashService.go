package transactionjpcbservice

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type CarWashService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdatePriority(workOrderSystemNumber, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse)
}
