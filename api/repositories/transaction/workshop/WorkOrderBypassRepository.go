package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderBypassRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse)
	Bypass(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshopentities.WorkOrderQualityControl, *exceptions.BaseErrorResponse)
}
