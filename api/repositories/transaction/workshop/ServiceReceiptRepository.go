package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"

	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceReceiptRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, id int, request transactionworkshoppayloads.ServiceReceiptSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
}
