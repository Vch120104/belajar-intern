package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AtpmClaimRegistrationRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse)
}
