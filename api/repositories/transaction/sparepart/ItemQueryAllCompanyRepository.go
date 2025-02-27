package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemQueryAllCompanyRepository interface {
	GetAllItemQueryAllCompany(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemQueryAllCompanyDownload(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]transactionsparepartpayloads.GetItemQueryAllCompanyDownloadResponse, *exceptions.BaseErrorResponse)
}
