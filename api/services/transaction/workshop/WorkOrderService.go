package transactionworkshopservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderService interface {
	WithTrx(Trxhandle *gorm.DB) WorkOrderService
	Save(transactionworkshoppayloads.WorkOrderRequest) (bool, error)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
