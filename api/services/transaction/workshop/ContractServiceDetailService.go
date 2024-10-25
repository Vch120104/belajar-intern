package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ContractServiceDetailService interface {
	GetAllDetail(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
