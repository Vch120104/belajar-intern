package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	WithTrx(trxHandle *gorm.DB) ItemRepository
	GetAllItem(filterCondition []utils.FilterCondition) ([]masteritempayloads.ItemLookup, error)
	GetAllItemLookup(map[string]string) ([]map[string]interface{}, error)
	GetItemById(Id int) (masteritempayloads.ItemResponse, error)
	GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, error)
	GetItemCode(string) ([]map[string]interface{}, error)
	SaveItem(masteritempayloads.ItemResponse) (bool, error)
	ChangeStatusItem(Id int) (bool, error)
}
