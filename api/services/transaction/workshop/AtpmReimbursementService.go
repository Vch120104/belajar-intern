package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type AtpmReimbursementService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	New(req transactionworkshoppayloads.AtpmReimbursementRequest) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse)
	Save(claimsysno int, req transactionworkshoppayloads.AtpmReimbursementUpdate) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse)
	Submit(claimsysno int) (bool, *exceptions.BaseErrorResponse)
}
