package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemImportRepository interface {
	GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)

	GetItemImportbyId(tx *gorm.DB, Id int) (any, *exceptionsss_test.BaseErrorResponse)

	SaveItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptionsss_test.BaseErrorResponse)
}
