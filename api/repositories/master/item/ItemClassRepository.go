package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassRepository interface {
	GetAllItemClass(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	GetItemClassById(tx *gorm.DB, Id int) (masteritempayloads.ItemClassResponse, error)
	SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, error)
	ChangeStatusItemClass(tx *gorm.DB, Id int) (bool, error)
}
