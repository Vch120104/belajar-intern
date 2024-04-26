package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemImportService interface {
	GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	GetItemImportbyId(Id int) (any, *exceptionsss_test.BaseErrorResponse)
	SaveItemImport(req masteritementities.ItemImport) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptionsss_test.BaseErrorResponse)
}
