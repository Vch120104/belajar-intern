package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemModelMappingService interface {
	CreateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptionsss_test.BaseErrorResponse)
	GetItemModelMappingByItemId(itemId int, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
}
