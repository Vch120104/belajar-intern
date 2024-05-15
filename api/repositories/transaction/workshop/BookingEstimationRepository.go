package transactionworkshoprepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BookingEstimationRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptionsss_test.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	Void(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	CloseOrder(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
}
