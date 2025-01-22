package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LicenseOwncerChangeRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetHistoryByChassisNumber(chassisNumber string, tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
