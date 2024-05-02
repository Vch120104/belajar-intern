package transactionworkshoprepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptionsss_test.BaseErrorResponse)
	Save(transactionworkshoppayloads.WorkOrderRequest) (bool, error)
	Submit(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	Void(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	CloseOrder(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
}
