package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassService interface {
	WithTrx(trxHandle *gorm.DB) ItemClassService
	GetAllItemClass(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, error)
	SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, error)
     ChangeStatusItemClass(Id int) (bool, error)
}
