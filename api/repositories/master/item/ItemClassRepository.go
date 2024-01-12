package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassRepository interface {
	WithTrx(trxHandle *gorm.DB) ItemClassRepository
	GetAllItemClass(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, error)
	SaveItemClass(request masteritempayloads.ItemClassResponse) (bool, error)
	ChangeStatusItemClass(Id int) (bool, error)
}
