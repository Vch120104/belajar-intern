package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemTypeRepository interface {
	GetAllItemType(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemTypeById(tx *gorm.DB, id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	GetItemTypeByCode(tx *gorm.DB, itemTypeCode string) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	CreateItemType(tx *gorm.DB, request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse)
	SaveItemType(tx *gorm.DB, id int, request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse)
	ChangeStatusItemType(tx *gorm.DB, id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	GetItemTypeDropDown(tx *gorm.DB) ([]masteritempayloads.ItemTypeDropDownResponse, *exceptions.BaseErrorResponse)
}
