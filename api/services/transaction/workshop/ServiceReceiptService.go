package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"

	"after-sales/api/utils"
)

type ServiceReceiptService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse)
	Save(id int, request transactionworkshoppayloads.ServiceReceiptSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
}
