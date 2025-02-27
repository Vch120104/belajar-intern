package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type ItemQueryAllCompanyService interface {
	GetAllItemQueryAllCompany(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemQueryAllCompanyDownload(filterCondition []utils.FilterCondition) (*excelize.File, *exceptions.BaseErrorResponse)
}
