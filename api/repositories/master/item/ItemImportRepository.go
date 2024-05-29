package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemImportRepository interface {
	GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)

	GetItemImportbyId(tx *gorm.DB, Id int) (any, *exceptions.BaseErrorResponse)

	SaveItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	UpdateItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
}
