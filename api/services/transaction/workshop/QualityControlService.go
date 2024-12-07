package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type QualityControlService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse)
	Qcpass(id int, iddet int) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
	Reorder(id int, iddet int, payload transactionworkshoppayloads.QualityControlReorder) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
}
