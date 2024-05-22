package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BookingEstimationService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	Void(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
}
