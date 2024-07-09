package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"

	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceReceiptRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, id int, request transactionworkshoppayloads.ServiceReceiptSaveRequest) (bool, *exceptions.BaseErrorResponse)
}
