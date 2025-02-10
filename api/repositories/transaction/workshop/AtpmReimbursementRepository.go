package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AtpmReimbursementRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	New(tx *gorm.DB, req transactionworkshoppayloads.AtpmReimbursementRequest) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, claimsysno int, req transactionworkshoppayloads.AtpmReimbursementUpdate) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, claimsysno int) (bool, *exceptions.BaseErrorResponse)
}
