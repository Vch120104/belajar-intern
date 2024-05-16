package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	NewStatus(tx *gorm.DB, request transactionworkshopentities.WorkOrderMasterStatus) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptionsss_test.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	Void(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	CloseOrder(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
}
