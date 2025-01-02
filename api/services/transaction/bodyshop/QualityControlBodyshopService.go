package transactionbodyshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	"after-sales/api/utils"
)

type QualityControlBodyshopService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse)
	Qcpass(id int, iddet int) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
	Reorder(id int, iddet int, payload transactionbodyshoppayloads.QualityControlReorder) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
}
