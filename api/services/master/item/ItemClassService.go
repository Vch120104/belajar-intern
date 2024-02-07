package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
)

type ItemClassService interface {
	GetAllItemClass(filterCondition []utils.FilterCondition) []map[string]interface{}
	GetItemClassById(Id int) masteritempayloads.ItemClassResponse
	SaveItemClass(req masteritempayloads.ItemClassResponse) bool
	ChangeStatusItemClass(Id int) bool
}
