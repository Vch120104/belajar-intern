package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassRepository interface {
	GetAllItemClass(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetItemClassDropDown(tx *gorm.DB) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse)
	GetItemClassById(tx *gorm.DB, Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse)
	GetItemClassByCode(tx *gorm.DB, itemClassCode string) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse)
	SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemClass(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
