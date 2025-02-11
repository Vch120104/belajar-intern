package transactionworkshoprepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PrintGatePassRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PrintById(tx *gorm.DB, id int) (map[string]string, *exceptions.BaseErrorResponse)
}
