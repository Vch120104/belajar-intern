package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemImportRepository interface {
	GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	GetItemImportbyId(tx *gorm.DB, Id int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemImport(tx *gorm.DB, req masteritempayloads.ItemImportUploadRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	GetItemImportbyItemIdandSupplierId(tx *gorm.DB, itemId int, supplierId int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
}
