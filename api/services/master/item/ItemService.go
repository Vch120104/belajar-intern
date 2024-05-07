package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
)

type ItemService interface {
	GetAllItem(filterCondition []utils.FilterCondition) []masteritempayloads.ItemLookup
	GetAllItemLookup(map[string]string) []map[string]interface{}
	GetItemById(Id int) masteritempayloads.ItemResponse
	GetItemWithMultiId(MultiIds []string) []masteritempayloads.ItemResponse
	GetItemCode(string) []map[string]interface{}
	SaveItem(masteritempayloads.ItemResponse) bool
	ChangeStatusItem(Id int) bool
}
