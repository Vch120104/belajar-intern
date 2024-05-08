package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemModelMappingRepository interface {
	GetItemModelMappingByItemId(tx *gorm.DB, itemId int, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	CreateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptionsss_test.BaseErrorResponse)
}
