package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemClassService interface {
	GetAllItemClass(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse)
	SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemClass(Id int) (bool, *exceptions.BaseErrorResponse)
}
