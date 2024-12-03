package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderBypassRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse)
	Bypass(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshoppayloads.WorkOrderBypassResponseDetail, *exceptions.BaseErrorResponse)
}
