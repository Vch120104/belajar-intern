package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemTypeService interface {
	GetAllItemType(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemTypeById(id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	GetItemTypeByCode(itemTypeCode string) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	CreateItemType(request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse)
	SaveItemType(id int, request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse)
	ChangeStatusItemType(id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse)
	GetItemTypeDropDown() ([]masteritempayloads.ItemTypeDropDownResponse, *exceptions.BaseErrorResponse)
}
