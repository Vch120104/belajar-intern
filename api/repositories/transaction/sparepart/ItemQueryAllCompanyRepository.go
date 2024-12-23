package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemQueryAllCompanyRepository interface {
	GetAllItemQueryAllCompany(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}