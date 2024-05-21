package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemModelMappingRepository interface {
	GetItemModelMappingByItemId(tx *gorm.DB, itemId int, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	CreateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse)
	UpdateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse)
}
